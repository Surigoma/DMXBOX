---
title: Art-Net Connection Guide
parent: Configuration and Hardware
nav_order: 3
---

# Art-Net Connection Guide

[Configuration and Hardware](settings.md) | [Configuration Reference](configuration.md)

This guide sends lighting data from DMXBOX to an Art-Net node, then through its DMX output to fixtures.
See [Lighting and Output Specifications](lighting-specifications.md) for packet and output details.

## Topology and Scope

```text
Browser / HTTP / TCP → DMXBOX → Ethernet → Art-Net node → DMX fixtures
```

DMXBOX sends one 512-channel universe to the configured destination and Port-Address.
Adding `dmx.groups` does not add universes.
Automatic node discovery, destination lists, and independent patches per universe are not implemented.
Received Art-Net data is not used as lighting input to DMXBOX.

## Prepare the Connection

- A PC and Art-Net node with Ethernet connectivity
- The node's IP address, subnet mask, and DMX output configuration procedure
- Fixture DMX start addresses and channel modes
- A group configuration already tested in DMXBOX

Start with one receiver on a wired network to make destination and routing issues easier to isolate.

## Example IP Settings

This example assumes a dedicated local network. Choose addresses that do not conflict with existing devices.

| Device | IP address | Subnet mask |
| --- | --- | --- |
| PC wired NIC | `192.168.50.10` | `255.255.255.0` |
| Art-Net node | `192.168.50.20` | `255.255.255.0` |

1. Connect the PC and node directly or through the same switch.
2. Configure both on the same subnet, as shown above.
3. Configure the node's intended DMX output port to receive Art-Net.
4. Set DMXBOX's `output.artnet.addr` to the node's IP address.

`addr` is the receiver address, not the PC's own address.
Do not use CIDR notation such as `192.168.50.20/24` or append a port as in `192.168.50.20:6454`.
Set the port separately using `port`.

## DMXBOX Configuration

This fragment is the value of `output`. Keep the complete configuration, including `input` and `dmx`.

```json
{
  "target": ["console", "artnet"],
  "artnet": {
    "addr": "192.168.50.20",
    "port": 6454,
    "net": 0,
    "subuni": 0,
    "universe": 1
  }
}
```

Keeping `console` enabled helps compare logs with physical output.
If retaining FTDI or OSC, preserve their required settings and target selections as well.
The default port is 6454. The current implementation uses the configured port for both the destination and the PC's UDP listener.
Normally leave it at 6454 to match the receiving node.

Restart if saving asks you to do so. A normal reload also recreates output resources,
so verify node output again after a change.

## Matching Universe Addresses

An IP subnet mask and an Art-Net Sub-Net are separate settings.
DMXBOX's `subuni` is the Art-Net Sub-Net; `universe` is the four-bit Universe within it.

```text
Port-Address = net × 256 + subuni × 16 + universe
```

| net | subuni | universe | Port-Address in the packet |
| --- | --- | --- | --- |
| 0 | 0 | 0 | 0 |
| 0 | 0 | 1 | 1 |
| 0 | 1 | 0 | 16 |
| 1 | 0 | 0 | 256 |

DMXBOX accepts `net: 0–127`, `subuni: 0–15`, and `universe: 0–15`.
The configuration example sends Port-Address 1. The row for zero helps interpret existing configurations.

The receiver's display of “Universe 1” can refer to packet address 0 or 1, depending on its numbering convention.
Check its manual for zero-based or one-based numbering rather than matching the displayed number alone.
If the node has separate address settings per DMX port, check the port actually connected to the fixtures.

## Multiple NICs and VPNs

The implementation enumerates local interfaces and binds to the first interface whose subnet contains the destination IP.
There is no configuration field for directly selecting a NIC name or source IP.

- Review wired and Wi-Fi connections to avoid putting both on the same subnet.
- Avoid ambiguous overlaps with VPN or virtual NIC subnets.
- OS-level reachability is not sufficient: initialization fails if no local subnet contains the destination.
- Sending through a router to another subnet is constrained by this source-selection behavior.

## First Output Test

1. Match the node's receiving Port-Address to DMXBOX.
2. Connect fixtures to that DMX output and match their channel modes and start addresses.
3. Check input restrictions such as Browser only, then operate a group from the browser.
4. Compare fixture response with the gauge during Fade In and Fade Out.
5. Verify Cut In and Cut Out, then test the same group from external control if required.

The gauge represents software output at the sender, not feedback from a fixture or confirmation of node reception.

## Diagnosing Missing Output

| Area | Check |
| --- | --- |
| DMXBOX | `artnet` is selected in `target`; the group contains devices |
| IP and NIC | Destination is the node's address; the PC has a matching subnet |
| UDP socket | Another lighting application is not using the same local IP and UDP port 6454 |
| Firewall | Network profile of the selected NIC and UDP permissions |
| Node | Art-Net input, Port-Address, and the actual DMX output port |
| Fixture | Channel mode, start address, and DMX cabling |

`Failed to contain network for interfaces.` means no matching local subnet was found.
For `Failed to create socket`, check availability and conflicts for the reported source IP and port.
`Drop` indicates a send error; logs do not guarantee cabling or receiver operation.

When using packet capture, filter on `udp.port == 6454` and inspect the destination IP,
ArtDMX Port-Address, and the values of channels you operated. Adjust the filter if using a custom port.
First establish whether packets leave the PC, then whether the node receives them, and finally inspect the DMX side.

## Shutdown and Signal Loss

DMXBOX does not explicitly send a blackout frame immediately before stopping.
Perform the required Fade Out and verify hold/blackout behavior on the node and fixtures.
Do not use network disconnection or PC shutdown as a substitute for a confirmed blackout operation.

## References

- [Official Art-Net glossary](https://art-net.org.uk/art-net-glossary/) — Port-Address, Net, Sub-Net, and Universe terminology
- [Operations and Troubleshooting](operations-guide.md) — logs and recovery

Art-Net™ Designed by and Copyright Artistic Licence.
