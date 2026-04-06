package templateservice

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"text/template"

	"github.com/minio/minio-go/v7"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/dynamic"
)

func (repo *Repository) renderManifest(rawYaml string, data interface{}) (string, error) {
	tmpl, err := template.New("k8s").Parse(rawYaml)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, data); err != nil {
		return "", err
	}

	return tpl.String(), nil
}

func (repo *Repository) ExecuteDeployment(ctx context.Context, templateName string, params DeployParams) error {
	files, err := repo.GetTemplatesConfigFile(ctx, templateName)
	if err != nil {
		return fmt.Errorf("gagal ambil config dari DB: %w", err)
	}

	log.Printf("Ditemukan %d file untuk template: %s", len(files), templateName)

	if len(files) == 0 {
		return fmt.Errorf("tidak ada file template untuk: %s", templateName)
	}

	bucketName := "templates"

	for i, f := range files {
		objectKey := fmt.Sprintf("%s/%s", strings.Split(f.Filename, "-")[0], f.Filename)

		log.Printf("[%d] Memproses: %s", i+1, objectKey)

		objInfo, err := repo.minioClient.StatObject(ctx, bucketName, objectKey, minio.StatObjectOptions{})
		if err != nil {
			errResp := minio.ToErrorResponse(err)
			if errResp.Code == "NoSuchKey" {
				return fmt.Errorf("file [%s] tidak ditemukan di bucket [%s]", objectKey, bucketName)
			}
			return fmt.Errorf("gagal stat object %s: %w", objectKey, err)
		}

		object, err := repo.minioClient.GetObject(ctx, bucketName, objectKey, minio.GetObjectOptions{})
		if err != nil {
			return fmt.Errorf("gagal ambil object %s: %w", objectKey, err)
		}

		buf := new(bytes.Buffer)
		buf.Grow(int(objInfo.Size))

		_, err = buf.ReadFrom(object)
		object.Close()

		if err != nil {
			return fmt.Errorf("gagal membaca konten %s: %w", objectKey, err)
		}

		renderedYaml, err := repo.renderManifest(buf.String(), params)
		if err != nil {
			return fmt.Errorf("gagal render file %s: %w", f.Filename, err)
		}

		fmt.Printf("\n--- [%d] APPLYING MANIFEST: %s ---\n", i+1, f.Filename)
		fmt.Println(renderedYaml)
		fmt.Println("------------------------------------------")

		err = repo.ApplyToK8s(ctx, renderedYaml, params.Namespace)
		if err != nil {
			log.Printf("failed to apply to k8s : %s", err.Error())
		}

	}

	log.Printf("Semua manifest untuk %s berhasil diproses!", templateName)
	return nil
}

func (repo *Repository) ApplyToK8s(ctx context.Context, renderedYaml string, namespace string) error {
	decUnstructured := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	obj := &unstructured.Unstructured{}
	_, gvk, err := decUnstructured.Decode([]byte(renderedYaml), nil, obj)
	if err != nil {
		return fmt.Errorf("failed to decode yaml: %w", err)
	}

	mapping, err := repo.mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return fmt.Errorf("failed to get rest mapping: %w", err)
	}

	var dr dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		dr = repo.dynClient.Resource(mapping.Resource).Namespace(namespace)
	} else {
		dr = repo.dynClient.Resource(mapping.Resource)
	}

	_, err = dr.Create(ctx, obj, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create resource %s: %w", obj.GetName(), err)
	}

	log.Printf("Successfully applied %s: %s", gvk.Kind, obj.GetName())
	return nil
}
