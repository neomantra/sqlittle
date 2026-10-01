#!/bin/bash
set -eu

# Needs an sqlite3 which reserves bytes at the end of each page (the macOS
# system sqlite3 does: 12 bytes). Checked in because it isn't reproducible
# everywhere.
rm -f reserved.sqlite
sqlite3 --batch reserved.sqlite <<HERE
CREATE TABLE big (id INTEGER PRIMARY KEY, s TEXT);
CREATE INDEX big_s ON big (s);
WITH RECURSIVE c(i) AS (SELECT 1 UNION ALL SELECT i+1 FROM c WHERE i < 12)
INSERT INTO big SELECT i, printf('%05d', i) || replace(hex(zeroblob(6000 + i*200)), '00', 'ab') FROM c;
HERE
