# wav

CLI tool for decoding WAV file headers.

## Example Usage

```
>> go run . samples/sine.wav

Size (bytes): 96036
Channels: 1
Sample rate (Hz): 48000
Byte rate (bytes/sec): 96000
Block align (bytes): 2
Bits per sample: 16
Data size (bytes): 96000
Duration: 1s
```

```
>> go run . samples/sine.wav --format json

{"size":96036,"channels":1,"sampleRate":48000,"byteRate":96000,"blockAlign":2,"bitsPerSample":16,"dataSize":96000,"duration":1}
```
