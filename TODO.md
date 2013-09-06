TODO
====

Lists
-----

 * RPUSH/RPOP
 * LPUSH/LPOP
 * Insert/delete

Checkpoints
-----------

Save shard as point-in-time snapshot.

1. Lock whole shard.
2. Rotate txlog.
3. Enter copy-on-write mode.
4. Unlock whole shard.
5. Copy whole keys to temp dataset on write. This gives ability to read changed data.
6. Write checkpoint to disk by iterating through all keys and dumping each one. (gob?)
7. When checkpoint has finished, lock whole shard.
8. Note last write in txlog and apply all writes back to in-memory copy.
9. Exit COW mode.
10. Unlock whole shard.
11. Delete in-memory copy.


### Notes

 * Don't allow >1 concurrent checkpoint per shard, obviously

Concurrency
-----------

To be defined :)

DONE
====

Maps
----

 * HPUT/HGET

Generic commutative operations
------------------------------

 * Change txlog format from a=b to OPERATION k v
 * func for each operation
 * Need Apply() func to switch on OPERATION
 * Introduces type compatibility error
