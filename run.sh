blade10:
	mkdir -p /tmp/memcockroach
	sudo mount -t tmpfs -o size=512M tmpfs /tmp/memcockroach/
	./cockroach init --stores=hdd="/tmp/memcockroach"
	./cockroach start --stores=hdd="/tmp/memcockroach" --addr=":26257" --gossip="astro9:26257,astro8:26257,astro6:26257,astro5:26257,astro4:26257"

blade9:
	mkdir -p /tmp/memcockroach
	sudo mount -t tmpfs -o size=512M tmpfs /tmp/memcockroach/
	./cockroach start --stores=hdd="/tmp/memcockroach" --addr=":26257" --gossip="astro10:26257,astro8:26257,astro6:26257,astro5:26257,astro4:26257"

blade9:
	mkdir -p /tmp/memcockroach
	sudo mount -t tmpfs -o size=512M tmpfs /tmp/memcockroach/
	./cockroach start --stores=hdd="/tmp/memcockroach" --addr=":26257" --gossip="astro9:26257,astro10:26257,astro6:26257,astro5:26257,astro4:26257"

blade9:
	mkdir -p /tmp/memcockroach
	sudo mount -t tmpfs -o size=512M tmpfs /tmp/memcockroach/
	./cockroach start --stores=hdd="/tmp/memcockroach" --addr=":26257" --gossip="astro9:26257,astro8:26257,astro10:26257,astro5:26257,astro4:26257"


blade9:
	mkdir -p /tmp/memcockroach
	sudo mount -t tmpfs -o size=512M tmpfs /tmp/memcockroach/
	./cockroach start --stores=hdd="/tmp/memcockroach" --addr=":26257" --gossip="astro9:26257,astro8:26257,astro6:26257,astro10:26257,astro4:26257"


blade9:
	mkdir -p /tmp/memcockroach
	sudo mount -t tmpfs -o size=512M tmpfs /tmp/memcockroach/
	./cockroach start --stores=hdd="/tmp/memcockroach" --addr=":26257" --gossip="astro9:26257,astro8:26257,astro6:26257,astro5:26257,astro10:26257"

