import os
import stat as statmod
from pathlib import Path, PurePosixPath, PureWindowsPath
import sys


class ManifestError(ValueError):
    pass


class ManifestFormatError(ManifestError):
    pass


class ManifestReadError(ManifestError):
    pass


class FilterError(ValueError):
    pass


class InvalidFilterError(FilterError):
    pass


class UnknownFixtureError(FilterError):
    pass


def display_text(value: str):
    return value.encode("utf-8", "surrogateescape").decode("utf-8", "replace")


def format_os_error(err: OSError):
    parts = []
    if getattr(err, "errno", None) is not None:
        parts.append(f"[Errno {err.errno}]")
    if getattr(err, "strerror", None):
        parts.append(display_text(str(err.strerror)))
    elif str(err):
        parts.append(display_text(str(err)))
    return " ".join(parts) or err.__class__.__name__


def configure_utf8_stdio():
    for stream in (sys.stdout, sys.stderr):
        if hasattr(stream, "reconfigure"):
            stream.reconfigure(encoding="utf-8")


def stat_workload_path(workload: str):
    path_arg = workload if os.name == "nt" else workload.encode("utf-8")
    return os.stat(path_arg)


def validate_fixture_name(name: str):
    reserved_device_names = {
        "con", "prn", "aux", "nul",
        "com1", "com2", "com3", "com4", "com5", "com6", "com7", "com8", "com9",
        "lpt1", "lpt2", "lpt3", "lpt4", "lpt5", "lpt6", "lpt7", "lpt8", "lpt9",
    }
    if len(name) > 200:
        raise ManifestFormatError(
            f"invalid compare-local fixture name {name!r}: must be 200 characters or fewer"
        )
    if name in (".", "..") or "/" in name or "\\" in name:
        raise ManifestFormatError(
            f"invalid compare-local fixture name {name!r}: must not contain path separators or dot-path segments"
        )
    if name.endswith("."):
        raise ManifestFormatError(
            f"invalid compare-local fixture name {name!r}: must not end with a dot"
        )
    device_stem = name.split(".", 1)[0].casefold()
    if device_stem in reserved_device_names:
        raise ManifestFormatError(
            f"invalid compare-local fixture name {name!r}: reserved device name"
        )
    if any(ch in name for ch in '<>:"|?*'):
        raise ManifestFormatError(
            f"invalid compare-local fixture name {name!r}: must not contain filesystem-special characters <>:\"|?*"
        )
    if "," in name:
        raise ManifestFormatError(
            f"invalid compare-local fixture name {name!r}: must not contain commas"
        )
    if name.casefold() == "all":
        raise ManifestFormatError(
            f"invalid compare-local fixture name {name!r}: reserved for the all-fixtures selector"
        )
    if name.startswith("--"):
        raise ManifestFormatError(
            f"invalid compare-local fixture name {name!r}: must not start with --"
        )


def resolve_workload_path(workload: str, manifest_dir: Path):
    posix_path = PurePosixPath(workload)
    windows_path = PureWindowsPath(workload)
    if workload.startswith("-"):
        raise ManifestFormatError(
            f"invalid compare-local workload path {workload!r}: must not start with -"
        )
    if not workload.lower().endswith(".xml"):
        raise ManifestFormatError(
            f"invalid compare-local workload path {workload!r}: must end with .xml"
        )
    if "\\" in workload:
        raise ManifestFormatError(
            f"invalid compare-local workload path {workload!r}: must use forward slashes instead of backslashes"
        )
    if posix_path.is_absolute() or windows_path.is_absolute():
        raise ManifestFormatError(
            f"invalid compare-local workload path {workload!r}: must not be absolute"
        )
    if windows_path.drive:
        raise ManifestFormatError(
            f"invalid compare-local workload path {workload!r}: must not include a Windows drive prefix"
        )
    if any(part == ".." for part in posix_path.parts) or any(part == ".." for part in windows_path.parts):
        raise ManifestFormatError(
            f"invalid compare-local workload path {workload!r}: must be repo-relative without '..' segments"
        )
    candidates = [workload]
    manifest_relative = manifest_dir / Path(workload)
    manifest_relative_str = str(manifest_relative)
    if manifest_relative_str not in candidates:
        candidates.append(manifest_relative_str)

    last_error = None
    for candidate in candidates:
        try:
            stat_result = stat_workload_path(candidate)
        except FileNotFoundError:
            continue
        except NotADirectoryError:
            continue
        except OSError as err:
            last_error = err
            continue
        if not statmod.S_ISREG(stat_result.st_mode):
            raise ManifestFormatError(
                f"invalid compare-local workload path {workload!r}: must refer to a file"
            )
        return workload if candidate == workload else manifest_relative_str

    if last_error is not None:
        raise ManifestFormatError(
            f"invalid compare-local workload path {workload!r}: {format_os_error(last_error)}"
        )
    raise ManifestFormatError(
        f"invalid compare-local workload path {workload!r}: does not exist"
    )


def read_manifest(manifest_path: str):
    fixtures = []
    seen_names = {}
    seen_names_folded = {}
    manifest_display = display_text(manifest_path)
    manifest_dir = Path(manifest_path).resolve().parent
    try:
        lines = Path(manifest_path).read_text(encoding="utf-8-sig").splitlines()
    except FileNotFoundError:
        raise ManifestReadError(f"compare-local manifest not found: {manifest_display}")
    except UnicodeDecodeError as err:
        raise ManifestReadError(f"unable to decode compare-local manifest {manifest_display}: {err}")
    except OSError as err:
        raise ManifestReadError(f"unable to read compare-local manifest {manifest_display}: {format_os_error(err)}")

    for line_no, raw_line in enumerate(lines, start=1):
        line = raw_line.strip()
        if not line or line.startswith("#"):
            continue
        fields = line.split()
        if len(fields) != 2:
            raise ManifestFormatError(
                f"invalid compare-local manifest line {line_no} in {manifest_display}: {line!r}"
            )
        name, workload = fields
        validate_fixture_name(name)
        workload = resolve_workload_path(workload, manifest_dir)
        if name in seen_names:
            raise ManifestFormatError(
                f"duplicate compare-local fixture name {name!r} on line {line_no} in {manifest_display}; first seen on line {seen_names[name]}"
            )
        folded_name = name.casefold()
        if folded_name in seen_names_folded:
            raise ManifestFormatError(
                f"case-insensitive duplicate compare-local fixture name {name!r} on line {line_no} in {manifest_display}; first seen as {seen_names_folded[folded_name]!r}"
            )
        seen_names[name] = line_no
        seen_names_folded[folded_name] = name
        fixtures.append({"name": name, "workload": workload})
    return fixtures


def parse_filter(raw_filter: str):
    items = []
    seen = set()
    for raw_item in raw_filter.split(","):
        item = raw_item.strip()
        folded = item.casefold()
        if not item or folded in seen:
            continue
        seen.add(folded)
        items.append(item)
    if len(items) == 1 and items[0].casefold() == "all":
        return []
    return items


def normalize_filter(raw_filter: str):
    selected = parse_filter(raw_filter)
    if not selected:
        return "all"
    return ",".join(selected)


def select_fixtures(fixtures, raw_filter: str):
    selected = parse_filter(raw_filter)
    if not selected:
        return fixtures
    fixture_map = {fixture["name"].casefold(): fixture for fixture in fixtures}
    return [fixture_map[name.casefold()] for name in selected if name.casefold() in fixture_map]


def format_filter_error(fixtures, err: FilterError):
    if isinstance(err, InvalidFilterError):
        return f"invalid compare-local filter: {err}"
    names = "".join(f"  - {fixture['name']}\n" for fixture in fixtures)
    return f"unknown compare-local fixture: {err}\nknown fixtures:\n{names}"


def validate_filter(fixtures, raw_filter: str):
    selected = parse_filter(raw_filter)
    if not selected:
        stripped = raw_filter.strip()
        non_empty_tokens = [item.strip().casefold() for item in raw_filter.split(",") if item.strip()]
        if stripped == "" or non_empty_tokens == ["all"]:
            return
        raise InvalidFilterError("filter did not include any fixture names")
    if any(name.casefold() == "all" for name in selected) and len(selected) > 1:
        raise InvalidFilterError("'all' cannot be combined with specific fixtures")
    known = {fixture["name"].casefold() for fixture in fixtures}
    for name in selected:
        if name.casefold() not in known:
            raise UnknownFixtureError(name)
