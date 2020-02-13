module.exports = {
    reporters: [
        'default',
        ['jest-junit',
            {
            suiteName: 'UI Unit Tests',
            outputDirectory: 'build/test-results',
            outputName: 'unit-tests.xml'
            }
        ]
    ]
}
