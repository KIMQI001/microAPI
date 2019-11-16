package testFuctions

// 构建虚拟密封过程，来进行硬件测试

// 使用方法：

import (
	"bytes"
	"context"
	"crypto/rand"
	"github.com/FC/FC/Dev/rpcClient/mySectorBuilder"
	"github.com/filecoin-project/go-filecoin/address"
	"github.com/filecoin-project/go-filecoin/proofs/sectorbuilder"
	"github.com/filecoin-project/go-filecoin/proofs/verification"
	"github.com/filecoin-project/go-filecoin/types"
	go_sectorbuilder "github.com/filecoin-project/go-sectorbuilder"
	"github.com/ipfs/go-cid"
	dag "github.com/ipfs/go-merkledag"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/ssa/interp/testdata/src/fmt"
	"io"
	"log"
	"time"
)
//MaxTimeToSealASector represents the maximum amount of time the test should
//wait for a sector to be sealed. Seal performance varies depending on the
//computer, so we need to select a value which works for slow (CircleCI OSX
//build containers) and fast (developer machines) alike.
const MaxTimeToSealASector = time.Second * 360

// MaxTimeToGenerateSectorPoSt represents the maximum amount of time the test
// should wait for a proof-of-spacetime to be generated for a sector.
const MaxTimeToGenerateSectorPoSt = time.Second * 360

func TestAddpieceAndPost()  error{
	miner:=""
	dir:=""
	// TODO：ZOE：1：
	//  ：每次执行给 lastId +1,即会让sectorID+1
	lastId:=0
	sb,err:=mySectorBuilder.GenSectorBuilder(miner,dir,dir, uint64(lastId))
	if err != nil {
		log.Fatalln("Generate SectorBuilder failed~~~!!!!")
	}

	inputBytes := RequireRandomBytes(types.TwoHundredFiftySixMiBSectorSize.Uint64())
	ref, size, reader, err := CreateAddPieceArgs(inputBytes)
	startAddPiece:=time.Now()
	sectorID, err := sb.AddPiece(context.Background(), ref, size, reader)
	if err!=nil{
		return errors.Wrap(err,"addPiece failed~~~!!!")
	}

	//count add time
	endAddPiece:=time.Now()
	fmt.Sprintf("addPiece use %s",endAddPiece.Sub(startAddPiece).String())
	timeout := time.After(MaxTimeToSealASector + MaxTimeToGenerateSectorPoSt)

	select {
	case val := <-sb.SectorSealResults():
		if val.SealingResult.SectorID!=sectorID{
			return errors.Wrap(err,"sectorID is not same~~~!!!")
		}

		// count seal time
		endSeal:=time.Now()
		fmt.Sprintf("sealSector use %s",endSeal.Sub(endAddPiece).String())

		minerAddr,err := address.NewFromString(miner)
		if err != nil {
			return errors.Wrap(err,"gen miner addr failed~~~!!!")
		}

		sres, serr := (&verification.RustVerifier{}).VerifySeal(verification.VerifySealRequest{
			CommD:      val.SealingResult.CommD,
			CommR:      val.SealingResult.CommR,
			CommRStar:  val.SealingResult.CommRStar,
			Proof:      val.SealingResult.Proof,
			ProverID:   sectorbuilder.AddressToProverID(minerAddr),
			SectorID:   val.SealingResult.SectorID,
			SectorSize: types.TwoHundredFiftySixMiBSectorSize,
		})
		if serr!=nil{
			return errors.Wrap(serr,"verify failed~~!!")
		}
		if sres.IsValid!=true{
			return errors.New("verify is not valid~~!!")
		}

		// count verify seal time
		endverifySeal:=time.Now()
		fmt.Sprintf("sealSector use %s",endverifySeal.Sub(endSeal).String())

		// TODO: This should be generates from some standard source of
		// entropy, e.g. the blockchain
		challengeSeed := types.PoStChallengeSeed{1, 2, 3}

		// TODO：ZOE：2：
		//  ：将此次的val.SealingResult.CommR,导出
		//   将函数在此截断,
		//   将go_sectorbuilder.NewSortedSectorInfo()，单独拿出，导入刚才的commr 1,2,3，4,进行测试

		sortedSectorInfo := go_sectorbuilder.NewSortedSectorInfo(go_sectorbuilder.SectorInfo{CommR: val.SealingResult.CommR})

		// generate a proof-of-spacetime
		gres, gerr := sb.GeneratePoSt(sectorbuilder.GeneratePoStRequest{
			SortedSectorInfo: sortedSectorInfo,
			ChallengeSeed:    challengeSeed,
		})
		if gerr!=nil{
			return errors.Wrap(serr,"generatePost failed~~!!")
		}

		// verify the proof-of-spacetime
		vres, verr := (&verification.RustVerifier{}).VerifyPoSt(verification.VerifyPoStRequest{
			ChallengeSeed:    challengeSeed,
			SortedSectorInfo: sortedSectorInfo,
			Faults:           []uint64{},
			Proof:            gres.Proof,
			SectorSize:       types.TwoHundredFiftySixMiBSectorSize,
		})
		if verr!=nil{
			return errors.Wrap(serr,"verify post failed~~!!")
		}
		if vres.IsValid!=true{
			return errors.New("verify Post is not valid~~!!")
		}

		// count Post time
		endPost:=time.Now()
		fmt.Sprintf("sealSector use %s",endPost.Sub(endverifySeal).String())

		return nil
	case <-timeout:
		return errors.New("timed out waiting for seal to complete")
	}
}



// CreateAddPieceArgs creates a PieceInfo for the provided bytes
func CreateAddPieceArgs(pieceData []byte) (cid.Cid, uint64, io.Reader, error) {
	data := dag.NewRawNode(pieceData)

	return data.Cid(), uint64(len(pieceData)), bytes.NewReader(pieceData), nil
}

// RequireRandomBytes produces n-number of bytes
func RequireRandomBytes( n uint64) []byte { // nolint: deadcode
	slice := make([]byte, n)

	if _, err := io.ReadFull(rand.Reader, slice);err!=nil{
		log.Fatalln("RequireRandomBytes failed~~!!!")
	}

	return slice
}
