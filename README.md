gokv
====

A stupid-simple and probably regrettable key-value store written in Go.

gokv stores every write in an append-only file of the format

    key=value\n
    ...

where *value* is a JSON object, array or value.

If it's not clear, *this is really not meant to be used in production systems*.
