let
  nixpkgs = import (builtins.fetchGit {
    url = "https://github.com/nixos/nixpkgs";
    ref = "release-20.03";
    rev = "5272327b81ed355bbed5659b8d303cf2979b6953";
  }) {};
in with nixpkgs;
buildGoModule rec {
  pname = "turbo-geth";
  version = "2020.07.01-alpha";

  src = fetchFromGitHub {
    owner = "ledgerwatch";
    repo = pname;
    rev = version;
    sha256 = "064arlvvvq15m0y1nf9a0mplrlfxsc0v3vmv8sl5262pknlvfyk3";
  };

  modSha256 = "1yclsimwskyqy9h7gdw60frqbm5fpp32h8xvsvli4nghkpnl175n";


subPackages = [
"cmd/geth"
];

  meta = with stdenv.lib; {
    homepage = "https://geth.ethereum.org/";
    description = "Official golang implementation of the Ethereum protocol";
    license = with licenses; [ lgpl3 gpl3 ];
    maintainers = with maintainers; [ adisbladis lionello xrelkd ];
  };
}
