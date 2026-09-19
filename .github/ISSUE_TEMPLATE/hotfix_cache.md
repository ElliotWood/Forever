---
name: Numbers from the beta client
about: Send DBCache.bin, a scrubbed DamageMeter.bin, or both
title: 'Beta client: '
labels: client-data
---

<!--
Thank you. This is the one source of Forever data that cannot be looked up
anywhere, because it only exists in the cache your own client downloads it into.

Drag the file into this box to attach it.
-->

## The file

<!--
DBCache.bin, from:

  World of Warcraft\_classic_beta_\Cache\ADB\enUS\

Copy it out with the game closed - the client holds it open while running. The
whole enUS folder is fine too, and slightly more useful: the per-table .tmp files
beside it say which tables have hotfixes without any guesswork.
-->

## What it contains, so you know what you are sending

<!--
Checked byte by byte, not assumed: NPC dialogue, item names and Blizzard's tuning
rows. No character name, no account, no realm, no Battle.net tag, no IP. It is the
same data everyone in your region is served.

The files get read and thrown away. They are never committed to this repository.
-->

## Optional, but it helps to know

- Roughly what you have been doing in the beta: <!-- raiding, dungeons, questing, sat in a city -->
- Is the `Spell<number>.tmp` in that folder larger than 4 bytes? <!-- yes / no / did not look -->

<!--
The client only caches rows it has actually needed, so what you have depends on
what you have done. That last question is the interesting one: if your Spell cache
is bigger than 4 bytes, your client knows about spell hotfixes mine does not, and
something in the sim is out of date.
-->

## Damage meter, if you are sending one

<!--
DamageMeter.bin, from World of Warcraft\_classic_beta_\Cache\ - the only record of what
the server actually paid out, since addons cannot read damage in Forever.

Getting one out is fiddly, in this order:

  1. Open the damage meter and leave it open - it records nothing while closed.
  2. Fight things.
  3. Log out to character select. The file is only written when the meter flushes,
     and it flushes on logout, which is why the folder looks empty while playing.
  4. Copy it out before logging back in. Logging in deletes it and starts again.

It carries character names, yours and everyone you grouped with, so run it through the
scrubber first: https://elliotwood.github.io/Forever/classic/scrub/

That page reads the file in your browser and never uploads it; the names are gone before
anything leaves your machine. Attach the file it hands back, not the original.
-->

## Anything that looked wrong in game

<!-- Optional. A tooltip that disagrees with this site, a cooldown that is not what
it says, an ability that does not behave the way the sim models it. -->
