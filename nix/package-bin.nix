{
  fetchurl,
  allowUnpinned ? false,
  lib,
  stdenvNoCC,
}:
let
  release = import ./release-assets.nix;
  system = stdenvNoCC.hostPlatform.system;
  asset = release.assets.${system};
  hash =
    if asset.hash != null then
      asset.hash
    else if allowUnpinned then
      lib.fakeHash
    else
      throw "wails3 hash not pinned yet for ${system}; publish the release and pin its SHA256 in nix/release-assets.nix";
in
stdenvNoCC.mkDerivation {
  pname = "wails3";
  version = release.version;

  src = fetchurl {
    url = "https://github.com/RigelBuild/wails/releases/download/${release.tag}/${asset.name}";
    inherit hash;
  };
  sourceRoot = ".";
  dontConfigure = true;
  dontBuild = true;
  installPhase = ''
    runHook preInstall
    install -Dm755 wails3 $out/bin/wails3
    runHook postInstall
  '';

  meta = {
    description = "Prebuilt wails3 CLI from the RigelBuild wails fork";
    homepage = "https://github.com/RigelBuild/wails";
    mainProgram = "wails3";
    platforms = builtins.attrNames release.assets;
  };
}
