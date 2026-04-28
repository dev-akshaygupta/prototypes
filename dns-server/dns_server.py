import socket
from dnslib import DNSRecord, QTYPE, RR, A
from dnslib.server import DNSServer, BaseResolver

ROOT_SERVER = "1.1.1.1"

LOCAL_RECORDS = {
    "www.dummy.test": "1.2.3.4",
    "www.myapp.local": "127.0.0.1"
}

class DNSResolver(BaseResolver):

    def resolve(self, request, handler):
        # Extract domain name from incoming DNS request
        domain = str(request.q.qname).rstrip(".")

        # Create DNS reply packet
        reply = request.reply()

        print("Resolving: ", domain)

        if domain in LOCAL_RECORDS:
            reply.add_answer(
                RR(domain, QTYPE.A, rdata=A(LOCAL_RECORDS[domain]), ttl=60)
            )
            return reply

        ip = resolve(domain)    # Call resolver
        if not ip:
            print("Resolution failed!")
            return reply
        
        # Add final resolved IP into response
        reply.add_answer(
            RR(domain, QTYPE.A, rdata=A(ip), ttl=60)
        )
        
        return reply

def query_dns(server, domain):
    # Create a DNS Question packet
    # Ex: What is the IP for `domain`?
    q = DNSRecord.question(domain)

    # Create UDP socket (DNS uses UDP port 53)
    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    sock.settimeout(3)  # prevent hanging forever, if no response from server

    # Send DNS packet to target DNS server
    sock.sendto(q.pack(), (server, 53))

    # Receive response packet
    data, _ = sock.recvfrom(512)

    return DNSRecord.parse(data)    # Comvert raw data into readable DNS structure

def resolve(domain):
    
    # Start resolvtion from root DNS server
    nameServer = ROOT_SERVER

    while True:
        # Ask current nameserver for domain IP
        response = query_dns(nameServer, domain)

        # If answer section exsit, we found final IP
        if response.rr:
            return str(response.rr[0].rdata)
        
        # Check for Additional Section
        # It contains next nameserver's IP
        if response.ar:
            nameServer = str(response.ar[0].rdata)
        else:
            # If Additional Section is missing
            # we only get nameServer domain name (NS record)
            # Ex: ns1.example.com
            # Resolve that to get it's IP
            nameServer = resolve(str(response.auth[0].rdata))

# Create resolver object
resolver = DNSResolver()

# Start DNS server locally on port 8053
server = DNSServer(resolver, port=8053, address="localhost")

print("DNS server running on localhost:8053...")
server.start()