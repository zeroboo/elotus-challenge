package elotuschallenge

func findLength(nums1 []int, nums2 []int) int {
	var maxSubLen int = 0
	len1 := len(nums1)
	len2 := len(nums2)
	//Iterate through the first array
	for i := 0; i < len1; i++ {
		//Start a new subarray from the current index
		for j := 0; j < len2; j++ {
			//Check if the sub array is valid
			if nums1[i] == nums2[j] {

				//Sub array started, now check how long it is
				k := 0
				for i+k < len1 && j+k < len2 && nums1[i+k] == nums2[j+k] {
					k++
				}
				//Update the maximum length found
				if k > maxSubLen {
					maxSubLen = k
				}
			}
		}
	}

	return maxSubLen
}
