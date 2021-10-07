package delivery

import (
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const imagePullSecretName = "spaceship-registry"

func (d *delivery) ensureImagePullSecrets(namespace string) ([]corev1.LocalObjectReference, error) {
	_, err := d.kubernetes.CoreV1().Secrets(namespace).Get(d.ctx, imagePullSecretName, metav1.GetOptions{})
	if err != nil {
		if !apierrors.IsNotFound(err) {
			d.log.WithError(err).Error("Unable to get existing image pull secret")
			return nil, err
		}

		err = d.createImagePullSecret(namespace)
		if err != nil {
			return nil, err
		}
	}

	return []corev1.LocalObjectReference{{Name: imagePullSecretName}}, nil
}

func (d *delivery) createImagePullSecret(namespace string) error {
	agentSecret, err := d.core.GetSecret()
	if err != nil {
		return err
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      imagePullSecretName,
			Namespace: namespace,
		},
		Data: map[string][]byte{
			corev1.DockerConfigJsonKey: []byte(agentSecret.DockerConfigJson),
		},
		Type: corev1.SecretTypeDockerConfigJson,
	}

	_, err = d.kubernetes.CoreV1().Secrets(namespace).Create(d.ctx, secret, metav1.CreateOptions{})
	if err != nil {
		d.log.WithError(err).Errorf("Unable to create image pull secret in %s namespace", namespace)

		return err
	}

	return nil
}
