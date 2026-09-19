// Receives the two beta client files the sim asks players for, so that sending one is a
// drag and a click rather than a GitHub account, an issue template and an attachment.
//
//   POST /upload?kind=dbcache        body: the raw DBCache.bin
//   POST /upload?kind=damagemeter    body: the raw scrubbed DamageMeter.bin
//
// Returns { ok: true, receipt } on success. The receipt is the object key, so a sender can
// quote it if they want to say which upload was theirs.
//
// This is a public endpoint that takes binary from strangers, so it is deliberately a bad
// place to put anything else:
//
//   Only two shapes are accepted. A DBCache.bin has to start with the XFTH magic. A
//   DamageMeter.bin has to contain the scrubber's placeholder, which means an unscrubbed
//   file cannot be accepted even by accident - the names have to be gone before the bytes
//   are worth sending.
//
//   Both are capped well above a real file and far below anything worth hosting.
//
//   A per-IP daily count, keyed by a salted hash so no address is stored. That and the
//   format check are what stop it being free file hosting, which is the realistic risk.
//
// Nothing here logs an IP or writes one into a key. The count is keyed by a hash.

// The origin check is not a security boundary - anything can set a header with curl, and
// this endpoint is public by design. It stops another site embedding a page that quietly
// posts a visitor's file here, which is the browser-shaped version of the risk. The real
// guards are the format check, the size cap and the daily count.
const ALLOWED_ORIGINS = new Set([
	'https://elliotwood.github.io',
	// Local development against a built copy of the site.
	'http://127.0.0.1:8129',
	'http://localhost:8129',
]);

const KINDS = {
	dbcache: {
		maxBytes: 16 * 1024 * 1024,
		// 'XFTH' - the hotfix cache header.
		check: bytes => bytes.length > 44 && bytes[0] === 0x58 && bytes[1] === 0x46 && bytes[2] === 0x54 && bytes[3] === 0x48,
		reason: 'not a hotfix cache: it should start with XFTH',
	},
	damagemeter: {
		maxBytes: 4 * 1024 * 1024,
		// The scrubber replaces every class-tagged name with Player<n>, so a scrubbed file
		// always carries at least one. An original never does.
		check: bytes => indexOfAscii(bytes, 'Player1') >= 0,
		reason: 'no scrubbed names found: send the file the scrub page hands back, not the original',
	},
};

const DAILY_PER_IP = 20;

function indexOfAscii(bytes, text) {
	const needle = [...text].map(c => c.charCodeAt(0));
	outer: for (let i = 0; i + needle.length <= bytes.length; i++) {
		for (let j = 0; j < needle.length; j++) if (bytes[i + j] !== needle[j]) continue outer;
		return i;
	}
	return -1;
}

function cors(origin, extra = {}) {
	return {
		'Access-Control-Allow-Origin': origin,
		'Access-Control-Allow-Methods': 'POST, OPTIONS',
		'Access-Control-Allow-Headers': 'Content-Type',
		'Access-Control-Max-Age': '86400',
		...extra,
	};
}

function json(origin, body, status = 200) {
	return new Response(JSON.stringify(body), {
		status,
		headers: cors(origin, { 'Content-Type': 'application/json' }),
	});
}

/** A day-scoped, salted hash of the address, so the counter never stores one. */
async function ipKey(request) {
	const address = request.headers.get('CF-Connecting-IP') || 'unknown';
	const day = new Date().toISOString().slice(0, 10);
	const digest = await crypto.subtle.digest('SHA-256', new TextEncoder().encode(`${day}:${address}`));
	return `count:${day}:${[...new Uint8Array(digest).slice(0, 8)].map(b => b.toString(16).padStart(2, '0')).join('')}`;
}

export default {
	async fetch(request, env) {
		const origin = request.headers.get('Origin') || '';
		const allowed = ALLOWED_ORIGINS.has(origin) ? origin : '';
		if (request.method === 'OPTIONS') return new Response(null, { status: 204, headers: cors(allowed) });

		const url = new URL(request.url);
		if (url.pathname !== '/upload') return json(allowed, { ok: false, error: 'not found' }, 404);
		if (request.method !== 'POST') return json(allowed, { ok: false, error: 'POST only' }, 405);
		if (!allowed) return json(allowed, { ok: false, error: 'wrong origin' }, 403);

		const kind = KINDS[url.searchParams.get('kind')];
		if (!kind) return json(allowed, { ok: false, error: 'unknown kind' }, 400);

		const declared = Number(request.headers.get('Content-Length') || 0);
		if (declared > kind.maxBytes) return json(allowed, { ok: false, error: 'file too large' }, 413);

		const key = await ipKey(request);
		const seen = Number((await env.UPLOAD_KV.get(key)) || 0);
		if (seen >= DAILY_PER_IP) return json(allowed, { ok: false, error: 'enough for today, thank you' }, 429);

		const bytes = new Uint8Array(await request.arrayBuffer());
		// Checked again after reading: Content-Length is whatever the client claimed.
		if (bytes.length > kind.maxBytes) return json(allowed, { ok: false, error: 'file too large' }, 413);
		if (!bytes.length) return json(allowed, { ok: false, error: 'empty file' }, 400);
		if (!kind.check(bytes)) return json(allowed, { ok: false, error: kind.reason }, 422);

		const name = url.searchParams.get('kind');
		const receipt = `${name}/${new Date().toISOString().slice(0, 10)}/${crypto.randomUUID()}.bin`;
		await env.UPLOADS.put(receipt, bytes, {
			httpMetadata: { contentType: 'application/octet-stream' },
			customMetadata: {
				kind: name,
				bytes: String(bytes.length),
				// Free text from the sender, capped. Never trusted, never rendered as HTML.
				note: (url.searchParams.get('note') || '').slice(0, 200),
			},
		});
		await env.UPLOAD_KV.put(key, String(seen + 1), { expirationTtl: 60 * 60 * 26 });

		return json(allowed, { ok: true, receipt });
	},
};
