package timeentry

const (
	phaseScript = `
		on change 
			if my value is not empty send phase%dActive to #hours then remove @disabled from next <select/> within #phases 
			else send phase%dInactive to #hours then send clearPhase to the next <select/> within #phases 
		end 
		on clearPhase set my value to empty add @disabled to me send phase%dInactive to #hours then send clearPhase to the next <select/> within #phases end
	`
	phaseCountScript = `
		on phase1Active from #hours or phase2Active from #hours or phase3Active from #hours or phase4Active from #hours or phase5Active from #hours or phase5Inactive from #hours
			set :count to 0
			set :allPhases to <select/> in #phases
			for p in :allPhases
				if p's value is not empty increment :count by 1
			end
			set my value to :count
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
