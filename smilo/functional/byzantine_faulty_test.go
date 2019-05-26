package functional_test

import (
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	tests "go-smilo/src/blockchain/regression"
	"go-smilo/src/blockchain/regression/src/container"
)

var _ = Describe("QFS-05: Byzantine Faulty", func() {

	Context("QFS-05-01: F faulty fullnodes", func() {
		const (
			numberOfNormal = 3
			numberOfFaulty = 1
		)
		var (
			vaultNetwork container.VaultNetwork
			blockchain   container.Blockchain
			err          error
		)
		BeforeEach(func() {
			vaultNetwork, err = container.NewDefaultVaultNetwork(dockerNetwork, numberOfNormal+numberOfFaulty)
			Expect(err).To(BeNil())
			Expect(vaultNetwork).ToNot(BeNil())
			Expect(vaultNetwork.Start()).To(BeNil())
			blockchain, err = container.NewDefaultSmiloBlockchainWithFaulty(dockerNetwork, vaultNetwork, numberOfNormal, numberOfFaulty)
			Expect(err).To(BeNil())
			Expect(blockchain).ToNot(BeNil())
			Expect(blockchain.Start(true)).To(BeNil())
		})

		AfterEach(func() {
			Expect(blockchain.Stop(true)).To(BeNil())
			blockchain.Finalize()
			Expect(vaultNetwork.Stop()).To(BeNil())
			vaultNetwork.Finalize()
		})

		It("Should generate blocks", func(done Done) {

			By("Wait for p2p connection", func() {
				tests.WaitFor(blockchain.Fullnodes(), func(geth container.Ethereum, wg *sync.WaitGroup) {
					Expect(geth.WaitForPeersConnected(numberOfNormal + numberOfFaulty - 1)).To(BeNil())
					wg.Done()
				})
			})

			By("Wait for blocks", func() {
				const targetBlockHeight = 3
				tests.WaitFor(blockchain.Fullnodes()[:1], func(geth container.Ethereum, wg *sync.WaitGroup) {
					Expect(geth.WaitForBlocks(targetBlockHeight)).To(BeNil())
					wg.Done()
				})
			})

			close(done)
		}, 60)
	})

	Context("QFS-05-01: F+1 faulty fullnodes", func() {
		const (
			numberOfNormal = 2
			numberOfFaulty = 2
		)
		var (
			vaultNetwork container.VaultNetwork
			blockchain   container.Blockchain
			err          error
		)
		BeforeEach(func() {
			vaultNetwork, err = container.NewDefaultVaultNetwork(dockerNetwork, numberOfNormal+numberOfFaulty)
			Expect(err).To(BeNil())
			Expect(vaultNetwork).ToNot(BeNil())
			Expect(vaultNetwork.Start()).To(BeNil())
			blockchain, err = container.NewDefaultSmiloBlockchainWithFaulty(dockerNetwork, vaultNetwork, numberOfNormal, numberOfFaulty)
			Expect(err).To(BeNil())
			Expect(blockchain).ToNot(BeNil())
			Expect(blockchain.Start(true)).To(BeNil())
		})

		AfterEach(func() {
			Expect(blockchain.Stop(true)).To(BeNil())
			blockchain.Finalize()
			Expect(vaultNetwork.Stop()).To(BeNil())
			vaultNetwork.Finalize()
		})

		It("Should not generate blocks", func(done Done) {
			By("Wait for p2p connection", func() {
				tests.WaitFor(blockchain.Fullnodes(), func(geth container.Ethereum, wg *sync.WaitGroup) {
					Expect(geth.WaitForPeersConnected(numberOfNormal + numberOfFaulty - 1)).To(BeNil())
					wg.Done()
				})
			})

			By("Wait for blocks", func() {
				// Only check normal fullnodes
				tests.WaitFor(blockchain.Fullnodes()[:2], func(geth container.Ethereum, wg *sync.WaitGroup) {
					Expect(geth.WaitForNoBlocks(0, time.Second*30)).To(BeNil())
					wg.Done()
				})
			})
			close(done)
		}, 60)
	})

})
