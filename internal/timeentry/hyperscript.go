package timeentry

const (
	phaseScript = `
		on change 
			if my value is not empty send phase%dActive to #employeeDetail then remove @disabled from next <select/> within #phases 
			else send phase%dInactive to #employeeDetail then send clearPhase to the next <select/> within #phases 
		end 
		on clearPhase set my value to empty add @disabled to me send phase%dInactive to #employeeDetail then send clearPhase to the next <select/> within #phases end
	`
	employeeHoursScript = `
		on phase%dActive from #employeeDetail remove @disabled from me end 
		on phase%dInactive from #employeeDetail set my value to empty then add @disabled to me
	`
)
