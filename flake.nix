{
  description = "RigelBuild wails fork: prebuilt wails3 CLI from GitHub Releases";

  inputs.nixpkgs.url = "github:nixos/nixpkgs/nixos-26.05";

  outputs =
    { nixpkgs, ... }:
    let
      systems = [
        "aarch64-darwin"
        "x86_64-linux"
        "aarch64-linux"
      ];
      release = import ./nix/release-assets.nix;
      forSystems = nixpkgs.lib.genAttrs systems;
      packageFor =
        system: allowUnpinned:
        nixpkgs.legacyPackages.${system}.callPackage ./nix/package-bin.nix { inherit allowUnpinned; };
      checkFor =
        system:
        let
          pkgs = nixpkgs.legacyPackages.x86_64-linux;
          package = packageFor system true;
          asset = release.assets.${system};
          buildInputs = package.nativeBuildInputs ++ package.buildInputs ++ package.propagatedBuildInputs;
          inputNames = map (input: input.name or (builtins.baseNameOf (toString input))) buildInputs;
        in
        assert builtins.isString package.drvPath;
        assert package.src.name == asset.name;
        assert package.sourceRoot == ".";
        assert package.dontConfigure;
        assert package.dontBuild;
        assert !(builtins.any (name: nixpkgs.lib.hasInfix "go" name) inputNames);
        assert !(builtins.any (name: nixpkgs.lib.hasInfix "buildGoModule" name) inputNames);
        pkgs.runCommand "wails3-fetch-only-evaluation-${system}" (
          nixpkgs.lib.optionalAttrs (asset.hash != null) { inherit (package) src; }
        ) ''
          touch "$out"
        '';
    in
    {
      packages = nixpkgs.lib.genAttrs (builtins.filter (system: release.assets.${system}.hash != null) systems) (
        system:
        {
          wails3 = packageFor system false;
          default = packageFor system false;
        }
      );

      checks.x86_64-linux = builtins.listToAttrs (
        map (system: {
          name = "wails3-${system}";
          value = checkFor system;
        }) systems
      );
    };
}
