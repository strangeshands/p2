#!/bin/bash

RESULTS_DIR="results"

# Create results directory if it doesn't exist
mkdir -p "$RESULTS_DIR"

# Find next available result number
i=1

while [ -f "$RESULTS_DIR/result-$i.txt" ]; do
    ((i++))
done

OUTPUT_FILE="$RESULTS_DIR/result-$i.txt"

echo "Running program..."
echo "Saving output to $OUTPUT_FILE"

go run . | tee "$OUTPUT_FILE"

echo
echo "Done. Output saved to $OUTPUT_FILE"