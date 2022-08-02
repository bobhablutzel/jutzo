<template>
  <div class="text-container">
    <carousel :value="entries" :numVisible="3" :numScroll="1" ref="carousel"
              :autoplayInterval="3000"
              :circular="entries.length > 3" class="custom-carousel">
      <template #header>
        <h2>{{header}}</h2>
      </template>
      <template #item="slotProps">
        <div class="blog-entry">
          <div class="blog-entry-frame">
            <div class="blog-entry-header">
              <router-link :to="{ name: 'LibraryBlogEntry', params: { id: slotProps.data.id }}">{{ slotProps.data.title }}</router-link>
            </div>
            <div class="blog-entry-teaser">{{slotProps.data.teaser}}</div>
            <div class="blog-entry-footer">Publication date: {{slotProps.data.publicationDate.split('T')[0]}}</div>
          </div>
        </div>
      </template>
    </carousel>
  </div>
</template>

<script>

import {retrieveNewest} from "@/services/blogService";

export default {
  name: "LibraryBlogNewest",
  data() {
    return {
      header: 'Loading blog entries from api.hablutzel.com...',
      entries: [],
      responsiveOptions: [
        {
          breakpoint: '1200px',
          numVisible: 3,
          numScroll: 3
        },
        {
          breakpoint: '800px',
          numVisible: 2,
          numScroll: 2
        },
        {
          breakpoint: '400px',
          numVisible: 1,
          numScroll: 1
        }
      ]
    }
  },
  mounted() {
    retrieveNewest((entries) => {
      this.entries = entries
      this.header = 'Newest Blog Entries'

      // The carousel misbehaves if the number of entries
      // is less than the numVisible. So in this case we
      // force there to be enough entries by adding duplicates
      if (this.entries.length < 3) {
        this.entries = [...entries, ...entries, ...entries, ...entries]
      }
    })
  }
}
</script>

<style scoped>

.blog-entry {
  width: 100%;
  height: 100%;
  padding: 0 3px 0 3px ;
}

.blog-entry-frame {
  border: 1px solid var(--surface-border);
  border-radius: 3px;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.blog-entry-header {
  background-color: var(--gray-200);
  text-align: center;
  padding: 2px 2px 2px 2px;
}

.blog-entry-header a {
  font-size: 1.25rem;
  font-style: normal;
  text-decoration: none;
  color: black;
  -webkit-text-fill-color: black;
}

.blog-entry-teaser {
  margin: 10px 10px 10px 10px;
  text-align: left;
  font-size: 0.9rem;
}
.blog-entry-footer {
  font-size: 0.5rem;
  text-align: left;
  margin-left: 5px;
  justify-content: flex-end;
  flex-direction: column;
  display: flex;
  height: 100%
}

</style>
