#!/bin/bash

echo "Updating mock files for tests."
mockSources=$(cat ./build/mocks/mock_sources)
for mockFile in ${mockSources}
do
    mockFileDirectory=${mockFile%/*}/mocks/
    mockFileName=${mockFile##*/}
    mockDestination=${mockFileDirectory}${mockFileName/.go/_mock.go}
    echo "Generating mock for ${mockFile} at ${mockDestination}"
    mockgen -source=${mockFile} -destination=${mockDestination} -package=mocks
done
