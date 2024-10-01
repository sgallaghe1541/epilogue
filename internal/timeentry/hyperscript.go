package timeentry

const (
	phaseScript = `
		on change 
			if my value is not empty send phase%dActive to #hours then remove @disabled from next <select/> within #phases 
			else send phase%dInactive to #hours then send clearPhase to the next <select/> within #phases 
		end 
		on clearPhase set my value to empty add @disabled to me send phase%dInactive to #hours then send clearPhase to the next <select/> within #phases end
		on empEquipAdded from #phases 
			if my value is not empty send phase%dActive to #hours
			else send phase%dInactive to #hours
		end
	`
	employeeHoursScript = `
		on phase%dActive from #hours remove @disabled from me end 
		on phase%dInactive from #hours set my value to empty then add @disabled to me end
		on jobSelected from #phases set my value to empty then add @disabled to me end
	`
	equipmentHoursScript = `
		on phase%dActive from #hours remove @disabled from me end 
		on phase%dInactive from #hours set my value to empty then add @disabled to me end
		on jobSelected from #phases set my value to empty then add @disabled to me end
	`
)
