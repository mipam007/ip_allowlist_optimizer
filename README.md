# ip_allowlist_optimizer

This tool processes an IP list from an Excel file and aggregates them into subnets.

## Features
- Parses IP addresses from an Excel sheet.
- Aggregates IPs into the smallest possible subnets.
- Generates output in a format suitable for allowlists.

## Requirements
- The Excel sheet must have a tab named `iplist`.

## Usage
Run the application with the desired Excel file:

```bash
./ip_allowlist_optimizer alowlist.xlsx

