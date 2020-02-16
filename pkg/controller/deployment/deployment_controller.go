package deployment

import (
	"os"

	"github.com/nlopes/slack"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_deployment")

func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileDeployment{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

func add(mgr manager.Manager, r reconcile.Reconciler) error {
	c, err := controller.New("deployment-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	slackToken := os.Getenv("SLACK_TOKEN")
	slackChannel := os.Getenv("SLACK_CHANNEL")
	slackClient := slack.New(slackToken)

	source := &source.Kind{Type: &appsv1.Deployment{}}
	handler := &handler.EnqueueRequestForObject{}
	predicate := predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			oldContainers := e.MetaOld.(*appsv1.Deployment).Spec.Template.Spec.Containers
			newContainers := e.MetaNew.(*appsv1.Deployment).Spec.Template.Spec.Containers

			for oldIndex := 0; oldIndex < len(oldContainers); oldIndex++ {
				oldContainer := oldContainers[oldIndex]
				for newIndex := 0; newIndex < len(newContainers); newIndex++ {
					newContainer := newContainers[newIndex]
					if oldContainer.Name == newContainer.Name && oldContainer.Image != newContainer.Image {
						messageText := createMessage(newContainer.Name, newContainer.Image, oldContainer.Image)
						slackClient.PostMessage(slackChannel, slack.MsgOptionText(messageText, true))
					}
				}
			}

			return true
		},
	}

	err = c.Watch(source, handler, predicate)
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileDeployment{}

type ReconcileDeployment struct {
	client client.Client
	scheme *runtime.Scheme
}

func (r *ReconcileDeployment) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	return reconcile.Result{}, nil
}

func createMessage(name string, newImage string, oldImage string) string {
	return "[ " + oldImage + " => " + newImage + " ]" + " " +
		"deployment image updated" + " " +
		"( " + "deployment: " + name + ", " + "container: " + name + " )"
}
