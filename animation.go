package main

// AnimationAction represents a function called for each step of an animation that returns false when the animation is over.
type AnimationAction = func(int) bool

// Animation stores the action and current step index of an animation.
type Animation struct {
	action  AnimationAction // the action to perform at each step
	andThen *Animation      // an animation to start when this one finishes
	step    int             // the number of times action has been invoked
}

// Animator stores and applies a set of animations.
type Animator struct {
	animations map[int]*Animation
	nextIndex  int
}

// NewAnimator creates a new animator.
func NewAnimator() Animator {
	return Animator{
		animations: make(map[int]*Animation),
	}
}

// Animate adds an animation to the animator.
func (a *Animator) Animate(animation Animation) {
	a.animations[a.nextIndex] = &animation
	a.nextIndex++
}

// Step applies all running animations and returns whether any were running.
func (a *Animator) Step() bool {
	anyRunning := false
	for i, animation := range a.animations {
		anyRunning = true
		stillRunning := animation.action(animation.step)
		animation.step++
		if !stillRunning {
			if animation.andThen != nil {
				a.animations[i] = animation.andThen
			} else {
				delete(a.animations, i)
			}
		}
	}
	return anyRunning
}
