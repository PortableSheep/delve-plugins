#!/usr/bin/env python3
"""
Simple Registry Generator
Scans plugin directories and generates registry.yml from plugin.json and releases
"""

import os
import json
import yaml
import hashlib
from datetime import datetime
from pathlib import Path

def calculate_checksum(file_path):
    """Calculate SHA256 checksum of a file"""
    sha256_hash = hashlib.sha256()
    with open(file_path, "rb") as f:
        for chunk in iter(lambda: f.read(4096), b""):
            sha256_hash.update(chunk)
    return sha256_hash.hexdigest()

def scan_plugins():
    """Scan all plugin directories for releases"""
    plugins = {}

    # Find all directories with plugin.json
    for item in os.listdir('.'):
        plugin_dir = Path(item)
        if not plugin_dir.is_dir() or item.startswith('.'):
            continue

        plugin_json = plugin_dir / 'plugin.json'
        if not plugin_json.exists():
            continue

        print(f"Processing plugin: {item}")

        # Read plugin metadata
        with open(plugin_json) as f:
            plugin_data = json.load(f)

        # Process plugin releases
        releases_dir = plugin_dir / 'releases'
        versions = []

        if releases_dir.exists():
            for version_dir in releases_dir.iterdir():
                if version_dir.is_dir():
                    version_data = process_version(version_dir, plugin_data)
                    if version_data:
                        versions.append(version_data)

        if versions:
            # Sort versions by semver (newest first)
            versions.sort(key=lambda x: x['version'], reverse=True)

            plugins[item] = {
                'name': plugin_data['info']['name'],
                'description': plugin_data['info']['description'],
                'author': plugin_data['info']['author'],
                'license': plugin_data['info']['license'],
                'repository': plugin_data['info'].get('repository', ''),
                'homepage': plugin_data['info'].get('homepage', ''),
                'tags': plugin_data['info'].get('tags', []),
                'category': plugin_data['info'].get('category', 'utilities'),
                'min_delve_version': plugin_data['info'].get('min_delve_version', 'v0.1.0'),
                'versions': versions
            }

    return plugins

def process_version(version_dir, plugin_data):
    """Process a single version directory"""
    version = version_dir.name
    plugin_id = plugin_data['info']['id']

    # Find platform binaries
    assets = {}
    platforms = ['darwin-amd64', 'darwin-arm64', 'linux-amd64', 'linux-arm64', 'windows-amd64']

    for platform in platforms:
        binary_name = f"{plugin_id}-{platform}"
        if platform == 'windows-amd64':
            binary_name += '.exe'

        binary_path = version_dir / binary_name
        if binary_path.exists():
            assets[platform] = {
                'url': f"{plugin_id}/releases/{version}/{binary_name}",
                'checksum': f"sha256:{calculate_checksum(binary_path)}",
                'size': binary_path.stat().st_size
            }

    # Skip if no binaries found
    if not assets:
        print(f"  No binaries found for {version}")
        return None

    # Build version entry
    version_data = {
        'version': version,
        'released': datetime.fromtimestamp(version_dir.stat().st_mtime).isoformat() + 'Z',
        'compatibility': plugin_data.get('compatibility', ['v0.1.0', 'v0.2.0']),
        'assets': assets
    }

    # Add frontend if exists
    frontend_file = version_dir / 'component.js'
    if frontend_file.exists():
        version_data['frontend'] = {
            'entry': 'component.js',
            'checksum': f"sha256:{calculate_checksum(frontend_file)}"
        }

    # Add plugin metadata reference
    version_data['plugin_metadata'] = {
        'url': f"{plugin_id}/plugin.json",
        'checksum': f"sha256:{calculate_checksum(Path('.') / plugin_id / 'plugin.json')}"
    }

    print(f"  Found version {version} with {len(assets)} platform(s)")
    return version_data

def generate_registry(plugins):
    """Generate complete registry.yml structure"""
    return {
        'registry': {
            'name': 'Official Delve Plugin Registry',
            'description': 'Central registry for Delve plugins',
            'maintainer': 'portablesheep',
            'url': 'https://github.com/PortableSheep/delve-registry',
            'api_version': 'v1',
            'last_updated': datetime.utcnow().isoformat() + 'Z'
        },
        'plugins': plugins,
        'channels': {
            'stable': {
                'description': 'Stable, production-ready releases',
                'include_prerelease': False,
                'default': True
            },
            'beta': {
                'description': 'Beta releases for testing',
                'include_prerelease': True
            },
            'development': {
                'description': 'Latest development builds',
                'include_prerelease': True
            }
        },
        'categories': {
            'data-tools': {
                'name': 'Data Tools',
                'description': 'Plugins for data manipulation, formatting, and analysis'
            },
            'development-tools': {
                'name': 'Development Tools',
                'description': 'Plugins for software development workflows'
            },
            'monitoring': {
                'name': 'Monitoring & Analytics',
                'description': 'Plugins for monitoring systems and analyzing metrics'
            },
            'utilities': {
                'name': 'Utilities',
                'description': 'General purpose utility plugins'
            }
        },
        'api': {
            'base_url': 'https://raw.githubusercontent.com/PortableSheep/delve-registry/main',
            'endpoints': {
                'registry_metadata': '/registry.yml',
                'plugin_list': '/api/plugins.json',
                'plugin_details': '/api/plugins/{plugin_id}',
                'download_asset': '/{plugin_id}/releases/{version}/{asset_name}',
                'plugin_metadata': '/{plugin_id}/plugin.json'
            }
        },
        'stats': {
            'total_plugins': len(plugins),
            'total_versions': sum(len(p['versions']) for p in plugins.values()),
            'most_popular': list(plugins.keys())
        }
    }

def generate_api_files(plugins):
    """Generate API endpoint files"""
    api_dir = Path('api')
    api_dir.mkdir(exist_ok=True)

    plugins_dir = api_dir / 'plugins'
    plugins_dir.mkdir(exist_ok=True)

    # Generate plugins.json (discovery endpoint)
    plugin_list = []
    for plugin_id, plugin_data in plugins.items():
        latest_version = plugin_data['versions'][0] if plugin_data['versions'] else None
        plugin_list.append({
            'id': plugin_id,
            'name': plugin_data['name'],
            'description': plugin_data['description'],
            'author': plugin_data['author'],
            'category': plugin_data['category'],
            'tags': plugin_data['tags'],
            'latest_version': latest_version['version'] if latest_version else None,
            'min_delve_version': plugin_data['min_delve_version']
        })

    with open(api_dir / 'plugins.json', 'w') as f:
        json.dump({
            'plugins': plugin_list,
            'total': len(plugin_list),
            'last_updated': datetime.utcnow().isoformat() + 'Z'
        }, f, indent=2)

    # Generate individual plugin detail files
    for plugin_id, plugin_data in plugins.items():
        with open(plugins_dir / plugin_id, 'w') as f:
            json.dump(plugin_data, f, indent=2)

def main():
    """Main execution"""
    print("🔍 Scanning for plugins...")
    plugins = scan_plugins()

    if not plugins:
        print("❌ No plugins found")
        return

    print(f"✅ Found {len(plugins)} plugins")

    print("📝 Generating registry.yml...")
    registry = generate_registry(plugins)

    with open('registry.yml', 'w') as f:
        yaml.dump(registry, f, default_flow_style=False, sort_keys=False)

    print("🔌 Generating API files...")
    generate_api_files(plugins)

    print("✅ Registry generation complete!")
    print(f"   - registry.yml: {len(plugins)} plugins")
    print("   - api/plugins.json: discovery endpoint")
    print("   - api/plugins/*: individual plugin details")

if __name__ == '__main__':
    main()
