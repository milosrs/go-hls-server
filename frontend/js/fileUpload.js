export const randomString = (length) => (Math.random() + 1).toString(36).substring(length);

export const createFileChunks = async (file) => {
    const ret = []
    const streamReader = file.stream().getReader()

    while(true) {
        const {value, done} = await streamReader.read()
        if(done) {
            break
        }

        ret.push(new Int8Array(value))
    }

    return ret
}