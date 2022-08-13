<template>
  <div class="panel-container">
    <panel-menu :model="items"/>
    <divider layout="vertical"/>
    <div style="display: block">
      <router-view/>
    </div>
  </div>
</template>

<script>
export default {
  name: "FrameworkHome",
  data() {
    return {
      items: [
        { label:'Introduction', to: {name: 'FrameworkIntroduction' }},
        { label:'Overview', to: {name: 'FrameworkOverview' }, items: [
            { label: 'The Five Disciplines', to: { name: 'FrameworkFiveDisciplines' }},
            { label: 'Organizational Alignment', to: { name: 'FrameworkOrganizationalAlignment' }},
            { label: 'Federation', to: { name: 'UnderConstruction' }},
          ]
        },
        {
          label: 'Business & Strategy', to: {name: 'DisciplineBusinessStrategy' }, items: [
            { label: 'Overview', to: { name: 'DisciplineBusinessStrategy' }},
            { label: 'Strategy and Tactics', to: { name: 'UnderConstruction' }},
            { label: 'Business Capabilities', to: { name: 'UnderConstruction' }},
            { label: 'Differentiation', to: { name: 'UnderConstruction' }},
            { label: 'Naming Conventions', to: { name: 'UnderConstruction' }},
            { label: 'Technology Capabilities', to: { name: 'UnderConstruction' }},
          ]
        },
        { label:'Portfolio Management', items: [
            { label: 'Overview', to: { name: 'DisciplinePortfolioManagement' }}
          ]
        },
        { label:'Implementation', items: [
            { label: 'Overview', to: { name: 'DisciplineImplementation' }}
          ]
        },
        { label:'Guidance', items: [
            { label: 'Overview', to: { name: 'DisciplineGuidance' }}
          ]
        },
        { label:'Oversight', items: [
          { label: 'Overview', to: { name: 'DisciplineOversight' }}
        ]
        },
      ]
    }
  }

}
</script>

<style scoped>

</style>