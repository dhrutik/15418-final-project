# 15418-final-project

TITLE: An implementation of lock-free B+ Trees

URL: https://github.com/dhrutik/15418-final-project/edit/main/README.md

SUMMARY:
We are going to construct a lock-free implementation of the canonical data structure upon which most databases are built, the B+ Tree. We are going to compare this algorithm's efficiency with the more traditional B+ Tree implementation that makes use of locks, specifically latch crabbing.

BACKGROUND:
The main application in the real world of B+ Trees is in the realm of database systems. They are often used for efficient storage and retrieval of data. B+ trees are widely employed in database indexing, and making them lock-free can help us curb some of the latency in accessing the structure, bettering both concurrency and parallelism of access in a multi-threaded environment.

For some background on the data structure itself, B+ trees are balanced tree data structures where each node contains a certain number of keys and pointers. The keys are sorted, allowing for efficient search, insertion, and deletion operations. The tree's structure ensures that the data is stored in a way that optimizes range queries and sequential access. The locking procedure that is most often used it commonly known as "latch crabbing." This is technique used in concurrent B+ tree implementations to reduce contention and improve parallelism over write-heavy operations, such as node splits or merges. Traditional, or naive lock implementations can lead to contention issues, especially in high-concurrency scenarios. The latch crabbing technique allows for the minimal subsection of the tree to be contained in the critical section, to decrease the amount of contention, and latency lost as a result of it. Specifically, instead of holding one single, large lock for the entire duration of the operation, latch crabbing involves temporarily releasing locks as the operation progresses. The concurrent thread moves through the tree, releasing locks on nodes that are no longer deemed "critical" to the current operation. This allows other threads to access those nodes concurrently.  To illustrate a small example, if the left-most leaf node is to be split, this change may not propagate all the way to the right-most leaf, so we can deem the right-most leaf to be non-critical, and release the latch on it, so other threads can potentially perform operations on this remaining portion.

THE CHALLENGE:


RESOURCES:
Because the most interesting/relevant component of this project is not the base implementaion of B+ Trees, but rather tha conversion to a lock-free implementation, we will start with an existing implementation of the tree data structure. As of now, we intend to work in GoLang, and pull our initial code from this repository: https://github.com/collinglass/bptree. However, as we explore further, if we find that Go is a poorly-structured language to reach our goals as intended, we may pivot to a C++ base implementation (at which point we will update this README).

As of now, the main resource that we will use to guide our project is this paper,
https://dl.acm.org/doi/10.14778/3402707.3402719 (PALM: parallel architecture-friendly latch-free modifications to B+ trees on many-core processors).

Jason Sewall, Jatin Chhugani, Changkyu Kim, Nadathur Satish, and Pradeep Dubey. 2011. PALM: parallel architecture-friendly latch-free modifications to B+ trees on many-core processors. Proc. VLDB Endow. 4, 11 (August 2011), 795–806. https://doi.org/10.14778/3402707.3402719

GOALS AND DELIVERABLES:

PLATFORM CHOICE:

SCHEDULE:
