package timeentry

const (
	phaseScript = `
		on change 
			if my value is not empty send phase%dActive to #hours 
				if next <select/> within #phases != null remove @disabled from next <select/> end
			else send phase%dInactive to #hours
				if next <select/> within #phases != null send clearPhase to the next <select/> end
		end 
		on clearPhase set my value to '' add @disabled to me then set @disabled to 'disabled' then send phase%dInactive to #hours
			if next <select/> within #phases != null send clearPhase to the next <select/> end
		end
	`
	phaseCountScript = `
		on phase1Active from #hours or phase2Active from #hours or phase3Active from #hours or phase4Active from #hours or phase5Active from #hours or phase5Inactive from #hours
			set :count to 0
			set :allPhases to <select/> in #phases
			for p in :allPhases
				if p's value is not '' increment :count by 1
			end
			set my value to :count
		end
	`
	employeeHoursScript = `
		on phase%dActive from #hours set @phase-selected to 'y' then if @employee-selected == 'y' remove @disabled from me end 
		on phase%dInactive from #hours set my value to null then add @disabled to me then set @disabled to 'disabled' end
		on jobSelected from #phases 
			set my value to null 
			add @disabled to me 
			set @disabled to 'disabled' 
			set @phase-selected to 'n'
		end
		on employeeSelected from #hours 
			if @phase-selected == 'y' and @employee-selected == 'y'
				remove @disabled from me
			end
		end
	`
	equipmentHoursScript = `
		on phase%dActive from #hours set @phase-selected to 'y' then if @equipment-selected == 'y' remove @disabled from me end 
		on phase%dInactive from #hours set my value to null then add @disabled to me then set @disabled to 'disabled' end
		on jobSelected from #phases 
			set my value to null
			add @disabled to me 
			set @disabled to 'disabled' 
			set @phase-selected to 'n'
		end
		on equipmentSelected from #hours 
			if @phase-selected == 'y' and @equipment-selected == 'y'
				remove @disabled from me
			end
		end
	`
	employeeSelectScript = `
		on change 
			if my value is not empty
				set inputs to [<input/> in closest <div/>, <select/> in closest <div/>]
				for i in inputs 
					set i's @employee-selected to 'y'
				end
			else 
				set inputs to [<input/> in closest <div/>, <select/> in closest <div/>]
				for i in inputs 
					set i's @employee-selected to 'n'
				end
			end
			send employeeSelected to #hours
		end 
	`
	equipmentSelectScript = `
		on change 
			if my value is not empty
				set inputs to <input/> in closest <div/>
				for i in inputs 
					set i's @equipment-selected to 'y'
				end
			end
			send equipmentSelected to #hours
		end 
	`
	jobCertifiedScript = `
		on load 
			if my value == 'Y'
				send certifiedJobSelected to #employeeDetail
			else
				send nonCertifiedJobSelected to #employeeDetail
			end
		end
	`
	employeeClassScript = `
		on certifiedJobSelected from #employeeDetail set @certified-selected to 'y' 
			if @employee-selected == 'y' remove @disabled from me end 
		end
		on nonCertifiedJobSelected from #employeeDetail set @certified-selected to 'n' set @disabled to 'disabled' end
		on employeeSelected from #hours 
			if @certified-selected == 'y' and @employee-selected == 'y' 
				remove @disabled from me
			else 
				set @disabled to 'disabled'
			end
		end
	`
)
