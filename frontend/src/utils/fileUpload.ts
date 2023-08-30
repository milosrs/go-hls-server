import { FileUploadWithPreview } from 'file-upload-with-preview'

const fileCookie = 'file-name'
const uploadedCookie = 'uploaded-length'
let upload: FileUploadWithPreview | null = null

export const randomString = (length: number) => (Math.random() + 1).toString(36).substring(length);

export const createFileUpload = (hash: string) => {
    upload = new FileUploadWithPreview(hash);
    return hash
}

export const createFileChunks = async (alreadyUploaded: string): Promise<Uint8Array[]> => {
    const uploadedLength = parseInt(alreadyUploaded)
    const ret: Uint8Array[] = []
    const streamReader = upload.cachedFileArray[0].stream().getReader()

    while(true) {
        const {value, done} = await streamReader.read()
        if(done) {
            break
        }

        ret.push(new Uint8Array(value))
    }

    return ret
}

const prefix = '/api/users/'

export const sendFileChunks = async (chunks: Uint8Array[]) => {
    const progress = document.getElementById('progress')
    const fileName: HTMLInputElement = document.getElementById('fileName') as HTMLInputElement

    for(let i = 0; i < chunks.length; i++) {
        const body = JSON.stringify({
            Name: fileName.value,
            Bytes: Array.from(chunks[i]),
        })

        const resp = await (i === 0 ? sendStartUpload(chunks, fileName.value) : sendUploadChunk(chunks[i], fileName.value))

        if(resp.status === 200) {
            const currentProgress = (i / chunks.length) * 100
            progress.style.width = `${currentProgress.toFixed(0)}%`
            progress.innerText = `${currentProgress.toFixed(0)}%`
        } else {
            alert(`error.... ${resp.statusText}`)
        }
    }
}

const sendStartUpload = (chunk: Uint8Array[], name: string): Promise<any> => {
    const url = prefix + 'startUpload'
    const body = JSON.stringify({
        Name: name,
        NumberOfChunks: chunk.length,
        Bytes: Array.from(chunk[0])
    })

    return fetch(url, {
        headers: {
            'Content-Type': 'application/json',
        },
        method: 'POST',
        body,
    })
}

const sendUploadChunk = (chunk: Uint8Array, name: string): Promise<any> => {
    const url = prefix + 'uploadChunk'
    const body = JSON.stringify({
        Name: name,
        Bytes: Array.from(chunk)
    })

    return fetch(url, {
        headers: {
            'Content-Type': 'application/json',
        },
        method: 'PATCH',
        body,
    })
}