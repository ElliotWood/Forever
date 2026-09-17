#!/usr/bin/python

# Diffs the wording of every talent between the Forever beta client and Classic Era, ignoring the
# numbers. The trait curves say what a talent's values are; they do not say what it applies to, and
# that is where Forever has been quietly widening things - Improved Seals went from "your Seal of
# Righteousness and Judgement of Righteousness" to "your Seals and Judgements", Moonfury from three
# named spells to "your Arcane and Nature spells".
#
#   tools/data_watch/talent_text_diff.py                 # every talent whose wording moved
#   tools/data_watch/talent_text_diff.py --class druid   # one class
#   tools/data_watch/talent_text_diff.py --unreviewed    # only ones no beta pass names
#
# --unreviewed is a hint, not a verdict: the passes fixed more than their prose names, so a talent
# can be absent from the docs and still be implemented correctly. Read the sim before believing it.
#
# Talent ids come from ui/core/talents/trees/*.json, which is also where a tree can carry another
# class's spell id - a talent whose Forever text is about Fire Blast on the druid tree is a
# collision, not a druid change.

import argparse
import glob
import io
import json
import os
import re
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
from spell_client import ERA, FOREVER, Client

REPO = os.path.join(os.path.dirname(os.path.abspath(__file__)), '..', '..')


def wording(text):
	"""The sentence with the numbers and the variable references taken out."""
	text = re.sub(r'\$[a-zA-Z0-9@]+\d*[a-zA-Z0-9]*', 'X', text or '')
	text = re.sub(r'\d+(\.\d+)?', 'N', text)
	return re.sub(r'\s+', ' ', text).strip().lower()


def reviewed_names():
	seen = ''
	for path in glob.glob(os.path.join(REPO, 'docs', 'beta-pass', '*.md')):
		seen += io.open(path, encoding='utf-8').read()
	checklist = os.path.join(REPO, 'docs', 'forever_beta_checklist.md')
	if os.path.exists(checklist):
		seen += io.open(checklist, encoding='utf-8').read()
	return seen


def main():
	parser = argparse.ArgumentParser(description='Talent wording, Forever against Era.')
	parser.add_argument('--class', dest='klass', help='one class tree')
	parser.add_argument('--unreviewed', action='store_true', help='only talents no beta pass names')
	args = parser.parse_args()

	f, e = Client(FOREVER), Client(ERA)
	docs = reviewed_names() if args.unreviewed else ''
	trees = sorted(glob.glob(os.path.join(REPO, 'ui', 'core', 'talents', 'trees', '*.json')))
	count = 0

	for path in trees:
		klass = os.path.basename(path).replace('.json', '')
		if args.klass and args.klass != klass:
			continue
		for tree in json.load(io.open(path, encoding='utf-8')):
			for talent in tree.get('talents', []):
				ids = talent.get('spellIds') or []
				if not ids:
					continue
				forever, era = f.spell(ids[0]), e.spell(ids[0])
				if not forever or not era:
					continue
				if wording(forever['description']) == wording(era['description']):
					continue
				if args.unreviewed and talent['name'] in docs:
					continue
				count += 1
				print(f"\n## {klass} / {talent['name']} ({ids[0]})")
				print('  era     :', ' '.join(era['description'].split()))
				print('  forever :', ' '.join(forever['description'].split()))

	print(f"\n{count} talents whose wording moved")


if __name__ == '__main__':
	main()
