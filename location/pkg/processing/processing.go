package processing

import (
	"time"

	"github.com/safe-distance/socium-infra/location/pkg/models"
	cluster "github.com/smira/go-point-clustering"
)

// PerformDBScan sends Pings to the cluster.DBScan algorithm and returns ping clusters
func PerformDBScan(pings []models.Ping, eps float64, minPoints int) []models.PingCluster {
	// Get timestamp
	timestamp := time.Now()
	// Build point list
	points := cluster.PointList{}
	for _, ping := range pings {
		points = append(points, cluster.Point{ping.Lat, ping.Lon})
	}
	// Find clusters
	clusters, _ := cluster.DBScan(points, eps, minPoints)
	pc := make([]models.PingCluster, len(clusters))
	for i, c := range clusters {
		pc[i] = models.PingCluster{Count: c.C, Points: c.Points, Timestamp: timestamp}
	}
	return pc
}
