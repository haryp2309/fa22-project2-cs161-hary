#! /bin/python3

import os
from pathlib import Path
import re

CLIENT_GO = "client.go"

client_path = Path(os.path.dirname(os.path.realpath(__file__)))/"client"
splitted_client_path = client_path.parent/"client_splitted"

os.rename(client_path, splitted_client_path)

files = os.listdir(splitted_client_path)
files.remove(CLIENT_GO)
files.insert(0, CLIENT_GO)

mergedContent = ""
for filename in files:
    full_filename = splitted_client_path/filename
    with open(full_filename, "r") as f:
        content = f.read()
        if filename != CLIENT_GO:
            content = re.sub(r".*import\s*\([^)]*\)(.*)", "", content)
            content = re.sub(r"^package client", "", content)
            content = content.strip()
            content = f"\n// FILE: {filename.upper()}\n\n"+content
        mergedContent += content
    os.rename(full_filename, str(full_filename)+".bak")

os.mkdir(client_path)
with open(client_path/CLIENT_GO, "w") as f:
    f.write(mergedContent)