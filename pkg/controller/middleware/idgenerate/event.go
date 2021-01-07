package idgenerate

import (
	"fmt"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
)

const (
	IdGeneratorEventReasonOnAdded = "IdGeneratorOnAdded"
	IdGeneratorEventReasonOnDelete = "IdGeneratorOnDelete"
	IdGeneratorEventReasonOnUpdate = "IdGeneratorOnUpdate"
)

func eventMessage(generate *crdv1alpha1.IdGenerate, eventType string)string{
	return fmt.Sprintf("idgenerator '%s:%s' %s", generate.GetName(), generate.GetNamespace(), eventType)
}
