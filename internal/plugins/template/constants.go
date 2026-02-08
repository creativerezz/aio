package template

// InputSentinel is used to temporarily replace {{input}} during template processing
// to prevent recursive variable resolution
const InputSentinel = "__AIO_INPUT_SENTINEL_TOKEN__"
