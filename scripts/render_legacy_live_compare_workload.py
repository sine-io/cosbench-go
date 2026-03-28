#!/usr/bin/env python3

import sys
import xml.etree.ElementTree as ET
from pathlib import Path


def parse_config(raw):
    items = []
    for part in raw.split(";"):
        part = part.strip()
        if not part:
            continue
        if "=" in part:
            key, value = part.split("=", 1)
        else:
            key, value = part, ""
        items.append((key.strip(), value.strip()))
    return items


def is_placeholder(value):
    return value.startswith("<") and value.endswith(">")


def render_config(raw, endpoint, access_key, secret_key, region, path_style):
    items = parse_config(raw)
    rendered = {}
    for key, value in items:
        if key == "accesskey":
            rendered[key] = access_key
            continue
        if key == "secretkey":
            rendered[key] = secret_key
            continue
        if key == "endpoint":
            rendered[key] = endpoint
            continue
        if key == "region":
            if region:
                rendered[key] = region
            elif not is_placeholder(value):
                rendered[key] = value
            continue
        if key == "path_style_access":
            if path_style:
                rendered[key] = path_style
            elif not is_placeholder(value):
                rendered[key] = value
            continue
        if is_placeholder(value):
            continue
        rendered[key] = value
    if "accesskey" not in rendered:
        rendered["accesskey"] = access_key
    if "secretkey" not in rendered:
        rendered["secretkey"] = secret_key
    if "endpoint" not in rendered:
        rendered["endpoint"] = endpoint
    if region and "region" not in rendered:
        rendered["region"] = region
    if path_style and "path_style_access" not in rendered:
        rendered["path_style_access"] = path_style
    return ";".join(f"{key}={value}" for key, value in rendered.items())


def main(argv):
    if len(argv) != 9:
        raise SystemExit(
            "usage: render_legacy_live_compare_workload.py <fixture> <output> <backend> <endpoint> <access_key> <secret_key> <region> <path_style>"
        )
    fixture_path = Path(argv[1])
    output_path = Path(argv[2])
    _backend = argv[3]
    endpoint = argv[4]
    access_key = argv[5]
    secret_key = argv[6]
    region = argv[7]
    path_style = argv[8]

    tree = ET.parse(fixture_path)
    root = tree.getroot()
    for storage in root.iter("storage"):
        storage.set(
            "config",
            render_config(
                storage.get("config", ""),
                endpoint=endpoint,
                access_key=access_key,
                secret_key=secret_key,
                region=region,
                path_style=path_style,
            ),
        )

    output_path.parent.mkdir(parents=True, exist_ok=True)
    tree.write(output_path, encoding="utf-8", xml_declaration=True)


if __name__ == "__main__":
    main(sys.argv)
