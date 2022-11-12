#! /bin/python3

import os
from pathlib import Path
import re

CLIENT_GO = "client.go"

client_path = Path(os.path.dirname(os.path.realpath(__file__)))/"client"
files = os.listdir(client_path)
files.remove(CLIENT_GO)
files.insert(0, CLIENT_GO)

mergedContent = ""
for filename in files:
    with open(client_path/filename, "r") as f:
        content = f.read()
        if filename != CLIENT_GO:
            content = re.sub(r".*import\s*\([^)]*\)(.*)", "", content)
            content = re.sub(r"^package client", "", content)
            content = content.strip()
            content = f"\n// FILE: {filename.upper()}\n\n"+content
        mergedContent += content

with open(client_path/"out.go", "w") as f:
    f.write(mergedContent)