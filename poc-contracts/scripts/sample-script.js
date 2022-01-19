// We require the Hardhat Runtime Environment explicitly here. This is optional
// but useful for running the script in a standalone fashion through `node <script>`.
//
// When running the script with `npx hardhat run <script>` you'll find the Hardhat
// Runtime Environment's members available in the global scope.
const hre = require("hardhat");

async function main() {
  // Hardhat always runs the compile task when running scripts with its command
  // line interface.
  //
  // If this script is run directly using `node` you may want to call compile
  // manually to make sure everything is compiled
  // await hre.run('compile');

  // We get the contract to deploy
  const TestErc20 = await hre.ethers.getContractFactory("WhiteBridgeCoinERC20");
  const testErc20 = await TestErc20.deploy(1500000000);

  await testErc20.deployed();

  console.log("Test erc 20 deployed to:", testErc20.address);

  const BaseledgerTest = await hre.ethers.getContractFactory("BaseledgerTest");
  const baseledgerTest = await BaseledgerTest.deploy(testErc20.address);

  await baseledgerTest.deployed();

  console.log("Baseledger test deployed to: ", baseledgerTest.address);
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
