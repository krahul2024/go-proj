Example usage of https server, and gRPC with different lang support

* Open, read, write, delete, move files.
    - Understand file paths, permissions, and errors.
    - Use io/ioutil, os, and bufio packages.
    - Uploading & downloading files via HTTP

* Learn to handle multipart/form-data for uploads.
    - Stream downloads instead of loading the full file into memory.
    - Handle large files efficiently with io.Copy or buffered reads/writes.
    - Streaming & processing files
    - Implement streaming for uploads/downloads, useful for huge files.
    - Process files while streaming (e.g., calculate checksum, compress/decompress, transform content on the fly).

* Advanced topics
    - gRPC streaming for file transfer.
    - Concurrent processing of file chunks (Go routines).
    - Chunked uploads/downloads for reliability and resume support.
    - Security considerations (permissions, sanitization, rate limiting).

* Move to C++
    - Reimplement your Go programs in C++ for performance.
    - Learn low-level file I/O (fstream, ifstream, ofstream) and memory mapping (mmap) for huge files.
    - Combine with networking (sockets) for custom file servers.
