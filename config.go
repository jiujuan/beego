package beego

import (
    "bufio"
    "bytes"
    "io"
    "os"
    "strconv"
    "strings"
    "sync"
    "unicode"
)

var (
    bComment = []byte{'#'}
    bEmpty = []byte{}
    bEqual = []byte{'='}
    bDQuote = []byte{'"'}
)

type Config struct {
    filename string
    comment map[int][]string
    data map[string]string
    offset map[string]int64
    sync.RWMutex
}

func LoadConfig(name string) (*Config, error) {
    file, err := os.Open(name)
    if err != nil {
        return nil, err
    }

    cfg := &Config{
        file.Name(),
        make(map[int][]string),
        make(map[string]string),
        make(map[string]int64),
        sync.RWMutex{},
    }
    cfg.Lock()
    defer cfg.Unlock()
    defer file.Close()

    var comment bytes.Buffer
    buf := bufio.NewReader(file)

    for nComment, off := 0, int64(1); ; {
        line, _, err := buf.ReadLine()
        if err == io.EOF {
            break
        }
        if bytes.Equal(line, bEmpty) {
            continue
        }

        off += int64(len(line))

        if bytes.HasPrefix(line, bComment) {
            line = bytes.TrimLeft(line, "#")
            line = bytes.TrimLeftFunc(line, unicode.IsSpace)
            comment.Write(line)
            comment.WriteByte('\n')
            continue
        }
        if comment.Len() != 0 {
            cfg.comment[nComment] = []string{comment.String()}
            comment.Reset()
            nComment++
        }
        val := bytes.SplitN(line, bEqual, 2)
        if bytes.HasPrefix(val[1], bDQuote) {
            val[1] = bytes.Trim(val[1], `"`)
        }
        key := strings.TrimSpace(string(val[0]))
        cfg.comment[nComment-1] = append(cfg.comment[nComment-1], key)
        cfg.data[key] = strings.TrimSpace(string(val[1]))
        cfg.offset[key] = off
    }
    return cfg, nil
}

func (c *Config) Bool(key string) (bool, error) {
    return strconv.ParseBool(c.data[key])
}

func (c *Config) Int(key string) (int, error) {
    return strconv.Atoi(c.data[key])
}

func (c *Config) Float(key string) (float64, error) {
    return strconv.ParseFloat(c.data[key], 64)
}

func (c *Config) String(key string) string {
    return c.data[key]
}