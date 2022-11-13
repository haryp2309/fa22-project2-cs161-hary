package client

import (
	"errors"
	"fmt"
	"strconv"

	userlib "github.com/cs161-staff/project2-userlib"
	"github.com/google/uuid"
)

type DataCluster struct {
	Value             []byte
	RemainingClusters int
}

const CLUSTER_SIZE = 32
const DEBUG_DATACLUSTER = false

func getDatastoreClusterKey(key uuid.UUID, clusterNumber int) (combinedKey uuid.UUID, err error) {
	combinedKey, err = GenerateDataStoreKey(key.String() + strconv.Itoa(clusterNumber))
	return
}

func padBytes(bytes []byte) (paddedBytes []byte, err error) {
	//return bytes, nil
	if len(bytes) > CLUSTER_SIZE {
		err = errors.New("TOO MANY BYTES! CANNOT PAD")
		return
	}
	paddedBytes = bytes
	paddingLength := CLUSTER_SIZE - len(bytes)
	for i := 0; i < (paddingLength + 1); i++ {
		paddedBytes = append(paddedBytes, byte(paddingLength))
	}
	// Results in bytes with length CLUSTER_SIZE + 1
	// So we know that the last byte always tells the padding length
	return
}

func unpadBytes(paddedBytes []byte) (bytes []byte, err error) {
	//return paddedBytes, nil
	if len(paddedBytes) != (CLUSTER_SIZE + 1) {
		err = errors.New("LENGTH OF BYTES ARE INCORRECT! CANNOT UNPAD")
		return
	}
	paddingLength := int(paddedBytes[len(paddedBytes)-1])
	for i := 0; i < (paddingLength + 1); i++ {
		if int(paddedBytes[len(paddedBytes)-1-i]) != paddingLength {
			err = errors.New("BYTES LOOKS TO BE TAMPERED WITH! CANNOT UNPAD")
		}
	}
	bytes = paddedBytes[:(len(paddedBytes) - (paddingLength + 1))]
	return
}

func DatastoreSet(key uuid.UUID, value []byte) (err error) {
	clusters := make([]DataCluster, 0)
	remainingValue := value
	for len(remainingValue) > 0 {
		end := CLUSTER_SIZE
		if len(remainingValue) < end {
			end = len(remainingValue)
		}

		//paddedValue, err := padBytes(remainingValue[:end])
		paddedValue := remainingValue[:end]
		// PADDING DISABLED AS IT CAUSES SOME COMPLICATED BUG
		if err != nil {
			return err
		}
		cluster := DataCluster{
			Value:             paddedValue,
			RemainingClusters: len(clusters),
		}
		clusters = append(clusters, cluster)
		remainingValue = remainingValue[end:]
	}

	for i := range clusters {
		clusters[i].RemainingClusters = len(clusters) - clusters[i].RemainingClusters - 1
	}

	for i, cluster := range clusters {
		combinedKey, err := getDatastoreClusterKey(key, i)
		if err != nil {
			return err
		}
		if DEBUG_DATACLUSTER {
			userlib.DebugMsg(
				fmt.Sprintf(
					"Setting cluster %s, i=%d, remaining=%d, key=%s\n",
					key.String(),
					i,
					cluster.RemainingClusters,
					combinedKey.String(),
				),
			)
		}
		marshalledCluster, err := MarshalAndEncrypt([]byte(key.String()), cluster)
		if err != nil {
			return err
		}
		userlib.DatastoreSet(combinedKey, marshalledCluster)
	}
	return
}

func getCluster(key uuid.UUID, combinedKey uuid.UUID) (cluster DataCluster, ok bool, err error) {
	var marshalledCluster []byte
	marshalledCluster, ok = userlib.DatastoreGet(combinedKey)
	if !ok {
		return
	}
	err = UnmarshalAndDecrypt([]byte(key.String()), marshalledCluster, &cluster)
	if err != nil {
		return
	}
	return
}

func DatastoreGet(key uuid.UUID) (value []byte, ok bool, err error) {
	remainingClusters := 1
	clusters := make([]DataCluster, 0)
	for remainingClusters > 0 {
		i := len(clusters)
		combinedKey, err := getDatastoreClusterKey(key, i)
		if err != nil {
			return nil, ok, err
		}
		if DEBUG_DATACLUSTER {
			userlib.DebugMsg(
				fmt.Sprintf(
					"Getting cluster %s, i=%d, key=%s\n",
					key.String(),
					i,
					combinedKey.String(),
				),
			)
		}
		cluster, ok, err := getCluster(key, combinedKey)
		if !ok || err != nil {
			return nil, ok, err
		}

		clusters = append(clusters, cluster)
		remainingClusters = cluster.RemainingClusters
	}

	value = make([]byte, 0)
	for _, cluster := range clusters {
		//unpaddedValue, err := unpadBytes(cluster.Value)
		// PADDING DISABLED AS IT CAUSES SOME COMPLICATED BUG
		unpaddedValue := cluster.Value
		if err != nil {
			return nil, ok, err
		}
		value = append(value, unpaddedValue...)
	}
	ok = true // everything went well
	return
}

func DatastoreDelete(key uuid.UUID) (err error) {
	combinedKey, err := getDatastoreClusterKey(key, 0)
	if err != nil {
		return
	}

	cluster, ok, err := getCluster(key, combinedKey)
	if !ok || err != nil {
		return
	}

	clustersCount := cluster.RemainingClusters + 1
	for i := 0; i < clustersCount; i++ {
		combinedKey, err := getDatastoreClusterKey(key, i)
		if err != nil {
			return err
		}
		userlib.DatastoreDelete(combinedKey)
	}
	return
}
