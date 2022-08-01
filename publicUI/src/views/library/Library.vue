<template>
  <div class="panel-container">
    <panel-menu :model="topics" class="navigation"/>
    <divider layout="vertical"/>
    <router-view/>
  </div>
</template>

<script>
import {retrieveNewest} from "@/services/blogService";

export default {
  name: "ArchitectureLibrary",
  data() {
    return {
        blog_entries: [],
        topics: [
          { label: 'Overview', to: {name: 'LibraryOverview'}},
          { label: 'Blog', to: {name: 'LibraryBlog'}, items: []},
          { label: 'Job Descriptions', to: {name: 'JobDescriptionList'}, items: [
              {label: 'Chief Architect', to: {name: 'ChiefArchitectJD'}},
            ]
          },
        ]
    }
  },
  mounted() {
    retrieveNewest((entries) => {
      entries.forEach( (entry) => {
        this.topics[1].items.push( { label: entry.title, to: { name: 'LibraryBlogEntry', params: { id: entry.id }}})
      })
    })
  }
}
</script>

<style scoped>

</style>

{ path: '/library', name: 'Library', component: Library,
children: [
{ path: '', name: 'FrameworkIntroduction', component: LibraryScaffold },
{ path: 'overview', name: 'FrameworkOverview', component: LibraryOverview },
{ path: 'jd', name: 'FrameworkFiveDisciplines', component: LibraryJobDescriptions },
]
},
