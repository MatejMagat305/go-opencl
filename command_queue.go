package opencl

import (
	"unsafe"
)

type commandQueueProperties uint32

type CommandQueue uint

func (cq CommandQueue) EnqueueBarrier() error {
	return enqueueBarrier(cq).getError()
}

func (cq CommandQueue) EnqueueNDRangeKernel(kernel Kernel, workDim uint, globalOffsets, globalWorkSizes, localWorkSizes []uint64) error {
	var offsets, gsizes, lsizes []clSize
	if len(globalOffsets) > 0 {
		offsets = make([]clSize, len(globalOffsets))
		for i := range globalOffsets {
			offsets[i] = clSize(globalOffsets[i])
		}
	}
	if len(globalWorkSizes) > 0 {
		gsizes = make([]clSize, len(globalWorkSizes))
		for i := range globalWorkSizes {
			gsizes[i] = clSize(globalWorkSizes[i])
		}
	}
	if len(localWorkSizes) > 0 {
		lsizes = make([]clSize, len(localWorkSizes))
		for i := range localWorkSizes {
			lsizes[i] = clSize(localWorkSizes[i])
		}
	}
	return enqueueNDRangeKernel(
		cq, kernel, workDim, offsets, gsizes, lsizes, 0, nil, nil,
	).getError()
}

func (cq CommandQueue) EnqueueReadBuffer(buffer Buffer, blockingRead bool, data *BufferData) error {
	return enqueueReadBuffer(
		cq, buffer, blockingRead, 0, clSize(data.DataSize), data.Pointer, 0, nil, nil,
	).getError()
}

func (cq CommandQueue) EnqueueReadImage(image Buffer, blockingRead bool, data *ImageData) error {
	origin := [3]clSize{clSize(data.Origin[0]), clSize(data.Origin[1]), clSize(data.Origin[2])}
	region := [3]clSize{clSize(data.Region[0]), clSize(data.Region[1]), clSize(data.Region[2])}
	return enqueueReadImage(
		cq, image, blockingRead, origin, region, 0, 0, data.Pointer, 0, nil, nil,
	).getError()
}

func (cq CommandQueue) EnqueueWriteBuffer(buffer Buffer, blockingWrite bool, data *BufferData) error {
	return enqueueReadBuffer(
		cq, buffer, blockingWrite, 0, clSize(data.DataSize), data.Pointer, 0, nil, nil,
	).getError()
}

type MapFlag uint32

const (
	MapFlagRead                  MapFlag = (1 << 0)
	MapFlagWrite                 MapFlag = (1 << 1)
	MapFlagWriteInvalidateRegion MapFlag = (1 << 2)
)

func (cq CommandQueue) EnqueueMapBuffer(buffer Buffer, blockingMap bool, flags []MapFlag, data *BufferData) error {
	var st clStatus

	mapFlags := MapFlag(0)
	for _, f := range flags {
		mapFlags |= f
	}
	ptr := enqueueMapBuffer(
		cq, buffer, blockingMap, mapFlags, 0, clSize(data.DataSize), 0, nil, nil, &st,
	)
	if st != CL_SUCCESS {
		return st.getError()
	}
	data.Pointer = unsafe.Pointer(ptr)

	return nil
}

func (cq CommandQueue) EnqueueMapImage(image Buffer, blockingMap bool, flags []MapFlag, data *ImageData) error {
	var st clStatus

	mapFlags := MapFlag(0)
	for _, f := range flags {
		mapFlags |= f
	}
	origin := [3]clSize{clSize(data.Origin[0]), clSize(data.Origin[1]), clSize(data.Origin[2])}
	region := [3]clSize{clSize(data.Region[0]), clSize(data.Region[1]), clSize(data.Region[2])}
	rowpitch, slicepitch := (*clSize)(&data.RowPitch), (*clSize)(&data.SlicePitch)
	ptr := enqueueMapImage(
		cq, image, blockingMap, mapFlags, origin, region, rowpitch, slicepitch, 0, nil, nil, &st,
	)
	if st != CL_SUCCESS {
		return st.getError()
	}
	data.Pointer = unsafe.Pointer(ptr)

	return nil
}

func (cq CommandQueue) EnqueueUnmapBuffer(buffer Buffer, data *BufferData) error {
	return enqueueUnmapMemObject(
		cq, buffer, data.Pointer, 0, nil, nil,
	).getError()
}

func (cq CommandQueue) Finish() error {
	return finishCommandQueue(cq).getError()
}

func (cq CommandQueue) Flush() error {
	return flushCommandQueue(cq).getError()
}

func (cq CommandQueue) Release() error {
	return releaseCommandQueue(cq).getError()
}

// GL

func (cq CommandQueue) EnqueueAcquireGLObjects(objects []Buffer) error {
	return enqueueAcquireGLObjects(
		cq, uint32(len(objects)), unsafe.Pointer(&objects[0]), 0, nil, nil,
	).getError()
}

func (cq CommandQueue) EnqueueReleaseGLObjects(objects []Buffer) error {
	return enqueueReleaseGLObjects(
		cq, uint32(len(objects)), unsafe.Pointer(&objects[0]), 0, nil, nil,
	).getError()
}
