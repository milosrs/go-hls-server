package view

type Color string

const Red Color = "red"
const Blue Color = "blue"
const Green Color = "green"

// Represents a component registry
type component string

const Alert component = "alert"
const FileUpload component = "file-upload"
const Progress component = "progress"
const ProgressLabels component = "progress-labels"
