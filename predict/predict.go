package predict

import (
	"fmt"
	_ "image/png"
	"os"

	"github.com/metno/go-tf-cnn-predict/stddev"
	tf "github.com/wamuir/graft/tensorflow"
	"github.com/wamuir/graft/tensorflow/op"
)

const ImageWidth = 128
const ImageHeight = 128

func init() {
	os.Setenv("TF_CPP_MIN_LOG_LEVEL", "1")
	os.Setenv("KMP_AFFINITY", "noverbose")
}

func argmax(arr []float32) int {
	var big float32 = -1.0
	var idx int = -1

	for i := 0; i < len(arr); i++ {
		if arr[i] >= big {
			big = arr[i]
			idx = i
		}
	}
	return idx
}

type Predictor struct {
	Model *tf.SavedModel
}

// declare method
func NewPredictor(modeldir string) (Predictor, error) {
	predictor := Predictor{}
	model, err := tf.LoadSavedModel(modeldir, []string{"serve"}, nil)
	if err != nil {
		fmt.Printf("Error loading saved model: %s\n", err.Error())
		return Predictor{}, err
	}
	predictor.Model = model
	fmt.Println("Mode loaded")
	return predictor, err
}

// Predict - Return index with highest probability
func (p Predictor) Predict(imagefile string, model *tf.SavedModel) (int, error) {
	float32arr, err := p.PredictionsArr(imagefile, model)
	if err != nil {
		return -1, err
	}
	cc := argmax(float32arr)
	return cc, err

}

// Predict - return class prediction and standard deviation
func (p Predictor) PredictWithDeviation(imagefile string, model *tf.SavedModel) (classpred int, deviation float32, err error) {
	float32arr, err := p.PredictionsArr(imagefile, model)
	if err != nil {
		return -1, -1, err
	}
	deviation = stddev.CalcDeviation(float32arr)

	classpred = argmax(float32arr)
	//fmt.Printf("dev: %0.3f, %+v", deviation, float32arr)
	return classpred, deviation, err

}

// PredictionsArr Return array of probabilities
func (p Predictor) PredictionsArr(imagefile string, model *tf.SavedModel) ([]float32, error) {

	tensor, err := makeTensorFromImage(imagefile)
	if err != nil {
		return []float32{}, err

	}

	result, runErr := model.Session.Run(
		map[tf.Output]*tf.Tensor{
			model.Graph.Operation("serving_default_conv2d_input").Output(0): tensor,

			//model.Graph.Operation("keep_prob").Output(0): keepProb,
		},
		[]tf.Output{
			model.Graph.Operation("StatefulPartitionedCall").Output(0),
		},
		nil,
	)

	if runErr != nil {
		return []float32{}, runErr

	}
	arr32 := result[0].Value().([][]float32)[0]

	argmax(arr32)
	//fmt.Printf("Result: %v\n", result[0].Value().([][]float32))
	return arr32, nil
}

// Convert the image in filename to a Tensor suitable as input to the cc-classifier model.
func makeTensorFromImage(filename string) (*tf.Tensor, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	// DecodeJpeg uses a scalar String-valued tensor as input.
	tensor, err := tf.NewTensor(string(bytes))
	if err != nil {
		return nil, err
	}
	// Construct a graph to normalize the image
	graph, input, output, err := constructGraphToNormalizeImage()
	if err != nil {
		return nil, err
	}

	// Execute that graph to normalize this one image
	session, err := tf.NewSession(graph, nil)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	normalized, err := session.Run(
		map[tf.Output]*tf.Tensor{input: tensor},
		[]tf.Output{output},
		nil)
	if err != nil {
		return nil, err
	}

	return normalized[0], nil

}

// The inception model takes as input the image described by a Tensor in a very
// specific normalized format (a particular image size, shape of the input tensor,
// normalized pixel values etc.).
//
// This function constructs a graph of TensorFlow operations which takes as
// input a JPEG-encoded string and returns a tensor suitable as input to the
// inception model.
func constructGraphToNormalizeImage() (graph *tf.Graph, input, output tf.Output, err error) {

	// - The model was trained after with images scaled to 128x128 pixels.
	// - The colors, represented as R, G, B in 1-byte each were converted to
	//   float using (value - Mean)/Scale.
	const (
		H, W  = ImageWidth, ImageHeight
		Mean  = float32(1.0)
		Scale = float32(255)
	)
	// - input is a String-Tensor, where the string the JPEG-encoded image.
	// - The inception model takes a 4D tensor of shape
	//   [BatchSize, Height, Width, Colors=3], where each pixel is
	//   represented as a triplet of floats
	// - Apply normalization on each pixel and use ExpandDims to make
	//   this single image be a "batch" of size 1 for ResizeBilinear.
	s := op.NewScope()
	input = op.Placeholder(s, tf.String)
	decode := op.DecodeJpeg(s, input, op.DecodeJpegChannels(3))

	output = op.Div(s,
		op.Sub(s,
			op.ResizeBilinear(s,
				op.ExpandDims(s,
					op.Cast(s, decode, tf.Float),
					op.Const(s.SubScope("make_batch"), int32(0))),
				op.Const(s.SubScope("size"), []int32{H, W})),
			op.Const(s.SubScope("mean"), Mean)),
		op.Const(s.SubScope("scale"), Scale))
	graph, err = s.Finalize()

	return graph, input, output, err
}
