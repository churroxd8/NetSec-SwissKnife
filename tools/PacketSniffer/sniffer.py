# Project: sniffer.py
# Author: churroxd8
# Description: A TCP/IP Sniffer that captures raw packets, decodes headers, and identifies flags (SYN, ACK, FIN)

import socket
import struct
import textwrap

def main():
    # Creates a raw socket (Root privileges needed!), and checks for privileges and OS
    try:
        conn = socket.socket(socket.AF_PACKET, socket.SOCK_RAW, socket.ntohs(3))
    except AttributeError:
        print("Error: This script is designed for Linux systems.")
        return
    except PermissionError:
        print("Error: You need to run this as root (sudo).")
        return
    print("[*] Sniffing packets. Waiting for packets...")

    while True:
        # raw_data from wire, addr from source info
        raw_data, addr = conn.recvfrom(65536)

        # checks ethernet frame (14 bytes from ethernet header)
        dest_mac, src_mac, eth_proto, data = ethernet_frame(raw_data)

        print('f\nEthernet Frame:')
        print(f'\tDestination: {dest_mac}, Source: {src_mac}, Protocol: {eth_proto}')

        # filters for protocol 8 (IPv4 traffic)
        if eth_proto == 8:
            version, header_length, ttl, proto, src, target, data = ipv4_packet(data)
            print(f'\tIPv4 Packer:')
            print(f'\t\tVersion: {version}, header Length: {header_length}, TTL: {ttl}')
            print(f'\t\tProtocol: {proto}, Source: {src}, Target: {target}')
            print('\t\tData:')
            print(format_multi_line('\t\t\t', data))

            if proto == 6:
                src_port, dest_port, sequence, acknowledgment, urg, ack, psh, rst, syn, fin, payload = tcp_segment(data)

                print(f'\t\tTCP Segment:')
                print(f'\t\t\tSource Port: {src_port}, Destination Port: {dest_port}')
                print(f'\t\t\tSequence: {sequence}, Acknowledgment: {acknowledgment}')
                print(f'\t\t\tFlags:')
                print(f'\t\t\tURG: {urg}, ACK: {ack}, PSH: {psh}, RST: {rst}, SYN: {syn}, FIN: {fin}')

                # Payload printing
                if len(payload) > 0:
                    print('\t\t\tTCP Payload:')
                    print(format_multi_line('\t\t\t\t', payload))

def ethernet_frame(data):
    # '!' = Network Data (Big Endian)
    # '6s' = 6 bytes (MAC address)
    # 'H' = Unsigned Short (2 bytes for Protocol)
    dest_mac, src_mac, proto = struct.unpack('! 6s 6s H', data[:14])
    return get_mac_addr(dest_mac), get_mac_addr(src_mac), socket.htons(proto), data[14:]

# Formats the MAC address
def get_mac_addr(bytes_addr):
    bytes_str = map('{:02x}'.format, bytes_addr)
    return ':'.join(bytes_str).upper()

# Unpacks IPv4 Packet
def ipv4_packet(data):
    version_header_length = data[0]
    # Some bit manipulation to get version and length
    version = version_header_length >> 4
    header_length = (version_header_length & 15) * 4

    # Unpack the standard IP header (TTL, Protocol, Src IP, Dest IP)
    ttl, proto, src, target = struct.unpack('! 8x B B 2x 4s 4s', data[:20])
    return version, header_length, ttl, proto, ipv4(src), ipv4(target), data[header_length:]

# Formats IP string
def ipv4(addr):
    return '.'.join(map(str, addr))

# Hex Dump
def format_multi_line(prefix, string, size=80):
    size -= len(prefix)
    if isinstance(string, bytes):
        string = ''.join(r'\x{:02x}'.format(byte) for byte in string)
        if size % 2:
            size -= 1
    return '\n'.join([prefix + line for line in textwrap.wrap(string, size)])

def tcp_segment(data):
    # '!' = Network (Big Endian)
    # 'H' = Unsigned Short (2 bytes) -> Ports
    # 'L' = Unsigned Long (4 bytes) -> Seq/Ack numbers
    # 'H' = Unsigned Short (2 bytes) -> Offset & Flags combined
    (src_port, dest_port, sequence, acknowlegment, offset_reserved_flags) = struct.unpack('! H H L L H', data[:14])

    # The 'offset' (header length) is the first 4 bits of the 5th variable.
    # We shift right by 12 to isolate them.
    offset = (offset_reserved_flags >> 12) * 4

    # The 'flags' are the last 6 bits
    # Bitwise AND (&) is used to mask out everything else
    flag_urg = (offset_reserved_flags & 32) >> 5
    flag_ack = (offset_reserved_flags & 16) >> 4
    flag_psh = (offset_reserved_flags & 8) >> 3
    flag_rst = (offset_reserved_flags & 4) >> 2
    flag_syn = (offset_reserved_flags & 2) >> 1
    flag_fin = offset_reserved_flags & 1

    return src_port, dest_port, sequence, acknowlegment, flag_urg, flag_ack, flag_psh, flag_rst, flag_syn, flag_fin, data[offset:]

if __name__ == '__main__':
    main()

    