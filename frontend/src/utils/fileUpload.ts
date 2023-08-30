import { FileUploadWithPreview } from 'file-upload-with-preview'

const fileCookie = 'file-name'
const uploadedCookie = 'uploaded-length'
let upload: FileUploadWithPreview | null = null

export const randomString = (length: number) => (Math.random() + 1).toString(36).substring(length);

export const createFileUpload = (hash: string) => {
    upload = new FileUploadWithPreview(hash);
    return hash
}

export const createFileChunks = async (alreadyUploaded: string): Promise<Int8Array[]> => {
    const uploadedLength = parseInt(alreadyUploaded)
    const ret: Int8Array[] = []
    const streamReader = upload.cachedFileArray[0].stream().getReader()

    while(true) {
        const {value, done} = await streamReader.read()
        if(done) {
            break
        }

        ret.push(new Int8Array(value))
    }

    console.log("RETURN: ", ret)
    return ret
}