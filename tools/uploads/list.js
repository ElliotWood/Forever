// Lists what is in the uploads bucket, because R2 has no object listing outside a binding -
// wrangler can get an object by key but cannot tell you which keys exist, and the dashboard
// is the only other place to look.
//
//   cd tools/uploads && npx wrangler dev --remote   # then: curl 127.0.0.1:8787
//
// Or run ./ls.sh, which does both and shuts the preview down again.
export default {
	async fetch(_request, env) {
		const listed = await env.UPLOADS.list({ include: ['customMetadata'] });
		return Response.json(
			listed.objects.map(object => ({
				key: object.key,
				size: object.size,
				uploaded: object.uploaded,
				note: object.customMetadata?.note || '',
			})),
		);
	},
};
