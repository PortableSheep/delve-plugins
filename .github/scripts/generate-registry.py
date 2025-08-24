#!/usr/bin/env python3
"""
Generate registry.yml and API files from plugin directories
"""

import json
import yaml
import hashlib
from datetime import datetime
from pathlib import Path

def scan_plugins():
    """Scan plugins directory and extract metadata"""
    plugins = {}
    plugins_dir = Path(".")

    for plugin_dir in plugins_dir.iterdir():
        if plugin_dir.is_dir() and plugin_dir.name != ".github":
            plugin_json = plugin_dir / "plugin.json"
            if plugin_json.exists():
                with open(plugin_json) as f:
                    plugin_data = json.load(f)
                plugins[plugin_dir.name] = process_plugin(plugin_dir, plugin_data)

    return plugins

def process_plugin(plugin_dir, plugin_data):
    """Process a single plugin and its releases"""
    releases = []
    releases_dir = plugin_dir / "releases"

    if releases_dir.exists():
        for version_dir in releases_dir.iterdir():
            if version_dir.is_dir():
                release = process_release(version_dir, plugin_data)
                if release:
                    releases.append(release)

    # Sort releases by version (newest first)
    releases.sort(key=lambda x: x["version"], reverse=True)

    return {
        "name": plugin_data["info"]["name"],
        "description": plugin_data["info"]["description"],
        "author": plugin_data["info"]["author"],
        "license": plugin_data["info"]["license"],
        "repository": plugin_data["info"].get("repository"),
        "homepage": plugin_data["info"].get("homepage"),
        "tags": plugin_data["info"].get("tags", []),
        "category": plugin_data["info"].get("category", "utilities"),
        "min_delve_version": plugin_data["compatibility"]["min_delve_version"],
        "versions": releases
    }

def process_release(version_dir, plugin_data):
    """Process a single release version"""
    assets = {}

    # Scan for platform binaries
    for platform in plugin_data["build"]["platforms"]:
        os_name, arch = platform.split("/")
        binary_pattern = plugin_data["build"]["assets"]["binary"]
        binary_name = binary_pattern.format(
            plugin_name=plugin_data["info"]["id"],
            os=os_name,
            arch=arch
        )

        if os_name == "windows":
            binary_name += ".exe"

        binary_path = version_dir / binary_name
        if binary_path.exists():
            assets[f"{os_name}-{arch}"] = {
                "url": f"{plugin_data['info']['id']}/releases/{version_dir.name}/{binary_name}",
                "checksum": f"sha256:{calculate_checksum(binary_path)}",
                "size": binary_path.stat().st_size
            }

    # Process frontend assets
    frontend_entry = plugin_data.get("frontend", {}).get("entry")
    if frontend_entry:
        frontend_path = version_dir / frontend_entry
        if frontend_path.exists():
            frontend_checksum = calculate_checksum(frontend_path)
        else:
            frontend_checksum = "missing"
    else:
        frontend_checksum = None

    if not assets:
        return None  # Skip releases with no assets

    return {
        "version": version_dir.name,
        "released": datetime.fromtimestamp(version_dir.stat().st_mtime).isoformat() + "Z",
        "compatibility": plugin_data["compatibility"]["supported_versions"],
        "assets": assets,
        "frontend": {
            "entry": frontend_entry,
            "checksum": f"sha256:{frontend_checksum}"
        } if frontend_checksum else None,
        "plugin_metadata": {
            "url": f"{plugin_data['info']['id']}/plugin.json",
            "checksum": f"sha256:{calculate_checksum(Path('.') / plugin_data['info']['id'] / 'plugin.json')}"
        }
    }

def calculate_checksum(file_path):
    """Calculate SHA256 checksum of a file"""
    sha256_hash = hashlib.sha256()
    with open(file_path, "rb") as f:
        for chunk in iter(lambda: f.read(4096), b""):
            sha256_hash.update(chunk)
    return sha256_hash.hexdigest()

def generate_registry(plugins):
    """Generate registry.yml"""
    registry = {
        "registry": {
            "name": "Official Delve Plugin Registry",
            "description": "Central registry for Delve plugins",
            "maintainer": "portablesheep",
            "url": "https://github.com/PortableSheep/delve-registry",
            "api_version": "v1",
            "last_updated": datetime.utcnow().isoformat() + "Z"
        },
        "plugins": plugins,
        "channels": {
            "stable": {
                "description": "Stable, production-ready releases",
                "include_prerelease": False,
                "default": True
            },
            "beta": {
                "description": "Beta releases for testing",
                "include_prerelease": True
            },
            "development": {
                "description": "Latest development builds",
                "include_prerelease": True
            }
        },
        "categories": {
            "data-tools": {
                "name": "Data Tools",
                "description": "Plugins for data manipulation, formatting, and analysis"
            },
            "development-tools": {
                "name": "Development Tools",
                "description": "Plugins for software development workflows"
            },
            "monitoring": {
                "name": "Monitoring & Analytics",
                "description": "Plugins for monitoring systems and analyzing metrics"
            },
            "utilities": {
                "name": "Utilities",
                "description": "General purpose utility plugins"
            }
        },
        "api": {
            "base_url": "https://raw.githubusercontent.com/PortableSheep/delve-registry/main",
            "endpoints": {
                "registry_metadata": "/registry.yml",
                "plugin_list": "/api/plugins.json",
                "plugin_details": "/api/plugins/{plugin_id}",
                "download_asset": "/{plugin_id}/releases/{version}/{asset_name}",
                "plugin_metadata": "/{plugin_id}/plugin.json"
            }
        },
        "stats": {
            "total_plugins": len(plugins),
            "total_versions": sum(len(p["versions"]) for p in plugins.values()),
            "most_popular": list(plugins.keys())
        }
    }

    return registry

def generate_api_files(plugins):
    """Generate API endpoint files"""
    # Generate plugins.json (plugin list)
    plugin_list = []
    for plugin_id, plugin_data in plugins.items():
        latest_version = plugin_data["versions"][0] if plugin_data["versions"] else None
        plugin_list.append({
            "id": plugin_id,
            "name": plugin_data["name"],
            "description": plugin_data["description"],
            "author": plugin_data["author"],
            "category": plugin_data["category"],
            "tags": plugin_data["tags"],
            "latest_version": latest_version["version"] if latest_version else None,
            "min_delve_version": plugin_data["min_delve_version"]
        })

    # Generate individual plugin detail files
    api_dir = Path("api")
    api_dir.mkdir(exist_ok=True)

    plugins_dir = api_dir / "plugins"
    plugins_dir.mkdir(exist_ok=True)

    # Write plugins.json
    with open(api_dir / "plugins.json", "w") as f:
        json.dump({
            "plugins": plugin_list,
            "total": len(plugin_list),
            "last_updated": datetime.utcnow().isoformat() + "Z"
        }, f, indent=2)

    # Write individual plugin files
    for plugin_id, plugin_data in plugins.items():
        with open(plugins_dir / plugin_id, "w") as f:
            json.dump(plugin_data, f, indent=2)

def main():
    """Main function"""
    print("Scanning plugins directory...")
    plugins = scan_plugins()

    print(f"Found {len(plugins)} plugins")
    for plugin_id, plugin_data in plugins.items():
        versions = len(plugin_data["versions"])
        print(f"  - {plugin_id}: {versions} versions")

    print("Generating registry.yml...")
    registry = generate_registry(plugins)

    with open("registry.yml", "w") as f:
        yaml.dump(registry, f, default_flow_style=False, sort_keys=False)

    print("Generating API files...")
    generate_api_files(plugins)

    print("Registry generation complete!")

if __name__ == "__main__":
    main()
