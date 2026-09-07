// Package batcher combines individual items into size- or time-triggered batches.
//
// A [Batcher] uses one serial worker and a bounded queue. [Batcher.Add] applies
// backpressure, while [Batcher.TryAdd] supports immediate rejection when full.
// [Batcher.Flush] waits for a queue boundary and [Batcher.Shutdown] drains accepted
// work. Handler errors remain visible through these operations; [WithOnResult]
// reports every completed batch. The package does not retry failed batches.
package batcher
