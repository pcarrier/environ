package main

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.starlark.net/starlark"
)

func s3func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	bucket := ""
	prefix := ""
	region := ""
	profile := ""
	endpoint := ""
	conservative := false
	if err := starlark.UnpackArgs(fn.Name(), args, kwargs, "bucket", &bucket, "region", &region, "prefix?", &prefix, "profile?", &profile, "endpoint?", &endpoint, "conservative?", &conservative); err != nil {
		return nil, err
	}
	if prefix == "" {
		prefix = "environ"
	}
	var cfg aws.Config
	var err error
	if profile != "" {
		cfg, err = config.LoadDefaultConfig(context.Background(), config.WithSharedConfigProfile(profile))
	} else {
		cfg, err = config.LoadDefaultConfig(context.Background())
	}
	if err != nil {
		return nil, err
	}
	cfg.Region = region
	var client *s3.Client
	if endpoint != "" {
		client = s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(endpoint)
			o.UsePathStyle = true
		})
	} else {
		client = s3.NewFromConfig(cfg)
	}
	return S3{
		client:       client,
		bucket:       bucket,
		prefix:       prefix,
		conservative: conservative,
	}, nil
}

type S3 struct {
	client       *s3.Client
	bucket       string
	prefix       string
	conservative bool
}

func (s S3) Get(key string) ([]byte, error) {
	resp, err := s.client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s.prefix + "/" + key),
	})
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}

func (s S3) Write(key string, value []byte) error {
	var err error
	fullKey := s.prefix + "/" + key
	if s.conservative {
		_, err = s.client.HeadObject(context.Background(), &s3.HeadObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(fullKey),
		})
		if err == nil {
			return nil // already exists
		}
	}
	input := &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fullKey),
		Body:   strings.NewReader(string(value)),
	}
	if !s.conservative {
		input.IfNoneMatch = aws.String("*")
	}
	_, err = s.client.PutObject(context.Background(), input)
	if err != nil && strings.Contains(err.Error(), "PreconditionFailed") {
		return nil
	}
	return err
}

func (s S3) String() string {
	return fmt.Sprintf("s3(%s, %s)", s.bucket, s.prefix)
}

func (s S3) Type() string {
	return "S3"
}

func (s S3) Freeze() {
}

func (s S3) Truth() starlark.Bool {
	return starlark.Bool(true)
}

func (s S3) Hash() (uint32, error) {
	return starlark.String(s.String()).Hash()
}
