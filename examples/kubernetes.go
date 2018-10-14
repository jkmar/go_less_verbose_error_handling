package kubernetes

func (DeploymentV1Beta1) Generate(genericParams map[string]interface{}) (runtime.Object, error) {
	// START OMIT
	args, err := getArgs(genericParams)
	if err != nil {
		return nil, err
	}

	envs, err := getEnvs(genericParams)
	if err != nil {
		return nil, err
	}

	params, err := getParams(genericParams)
	if err != nil {
		return nil, err
	}

	name, err := getName(params)
	if err != nil {
		return nil, err
	}

	labels, err := getLabels(params, name)
	if err != nil {
		return nil, err
	}
	// END OMIT

	count, err := strconv.Atoi(params["replicas"])
	if err != nil {
		return nil, err
	}

	podSpec, err := makePodSpec(params, name)
	if err != nil {
		return nil, err
	}

	imagePullPolicy := v1.PullPolicy(params["image-pull-policy"])
	if err = updatePodContainers(params, args, envs, imagePullPolicy, podSpec); err != nil {
		return nil, err
	}

	if err := updatePodPorts(params, podSpec); err != nil {
		return nil, err
	}

	count32 := int32(count)
	deployment := extensionsv1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
		Spec: extensionsv1beta1.DeploymentSpec{
			Replicas: &count32,
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: *podSpec,
			},
		},
	}
	return &deployment, nil
}

